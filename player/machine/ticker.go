package machine

import (
	"errors"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/song"
	"github.com/heucuva/optional"
)

type ticker struct {
	settings tickerSettings
	current  struct {
		tick  int
		row   index.Row
		order index.Order
	}
	next struct {
		row   optional.Value[tickerRowBreak]
		order optional.Value[index.Order]
	}
	songLoop struct {
		current int
		detect  map[index.Order]struct{}
	}
}

type tickerRowBreak struct {
	row        index.Row
	breakOrder bool
}

type tickerSettings struct {
	InitialOrder          index.Order
	InitialRow            index.Row
	SongLoopStartingOrder index.Order
	SongLoopCount         int
	PlayUntilOrder        optional.Value[index.Order]
	PlayUntilRow          optional.Value[index.Row]
}

func initTick[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](t *ticker, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], settings tickerSettings) error {
	t.settings = settings
	t.current.tick = 0
	t.current.row = 0
	t.current.order = 0
	t.next.row.Set(tickerRowBreak{
		row:        settings.InitialRow,
		breakOrder: false,
	})
	t.next.order.Set(settings.InitialOrder)

	nextRow, nextOrder, err := advanceRowOrder(t, m)
	if err != nil {
		return err
	}

	if row, set := nextRow.Get(); set {
		t.current.row = row
	}
	if order, set := nextOrder.Get(); set {
		t.current.order = order
	}

	if t.songLoop.detect == nil {
		t.songLoop.detect = make(map[index.Order]struct{})
	}

	t.songLoop.detect[t.current.order] = struct{}{}

	m.us.SetTracingTick(t.current.order, t.current.row, t.current.tick)

	if err := m.onOrderStart(); err != nil {
		return err
	}

	return m.onRowStart()
}

func runTick[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](t *ticker, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	if err := m.onTick(); err != nil {
		return err
	}

	tick := t.current.tick + 1
	row := t.current.row
	order := t.current.order
	rowAdvanced := false
	orderAdvanced := false
	done := false
	if tick >= m.tempo {
		tick = 0
		rowAdvanced = true
	}

	if rowAdvanced {
		if err := m.onRowEnd(); err != nil {
			return err
		}

		nextRow, nextOrder, err := advanceRowOrder(t, m)
		if err != nil {
			if !errors.Is(err, song.ErrStopSong) {
				return err
			}

			done = true
		}

		if r, set := nextRow.Get(); set {
			row = r
		}

		if o, set := nextOrder.Get(); set {
			order = o
			orderAdvanced = true

			if err := m.onOrderEnd(); err != nil {
				return err
			}
		}
	}

	traceValueChangeWithComment(m, "tick", t.current.tick, tick, "runTick")
	traceValueChangeWithComment(m, "row", t.current.row, row, "runTick")
	traceValueChangeWithComment(m, "order", t.current.order, order, "runTick")
	t.current.tick = tick
	t.current.row = row
	t.current.order = order
	m.us.SetTracingTick(t.current.order, t.current.row, t.current.tick)

	if !done {
		o, oset := t.settings.PlayUntilOrder.Get()
		r, rset := t.settings.PlayUntilRow.Get()
		if oset || rset {
			orderMatch := true
			if oset {
				orderMatch = (o == t.current.order)
			}

			rowMatch := true
			if rset {
				rowMatch = (r == t.current.row)
			}

			done = orderMatch && rowMatch
		}
	}

	if done {
		return song.ErrStopSong
	}

	if orderAdvanced {
		if err := m.onOrderStart(); err != nil {
			return err
		}
	}
	if rowAdvanced {
		if err := m.onRowStart(); err != nil {
			return err
		}
	}
	return nil
}

func advanceRowOrder[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](t *ticker, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) (optional.Value[index.Row], optional.Value[index.Order], error) {
	row := int(t.current.row)
	rowUpdated := false
	order := int(t.current.order)
	orderUpdated := false

	desiredRow, rowSet := t.next.row.Get()
	desiredOrder, orderSet := t.next.order.Get()

	if rowSet && orderSet {
		row = int(desiredRow.row)
		rowUpdated = true
		order = int(desiredOrder)
		orderUpdated = true
	} else if rowSet {
		row = int(desiredRow.row)
		rowUpdated = true
		if desiredRow.breakOrder {
			order++
			orderUpdated = true
		}
	} else if orderSet {
		order = int(desiredOrder)
		orderUpdated = true
		row = 0
		rowUpdated = true
	} else {
		row++
		rowUpdated = true
	}

	t.next.row.Reset()
	t.next.order.Reset()

	orderScanMax := len(m.songData.GetOrderList())
	orderScanIter := 0
	forceLoopDetect := false
orderScan:
	if orderScanIter >= orderScanMax {
		order = int(t.settings.SongLoopStartingOrder)
		orderUpdated = true
		forceLoopDetect = true
	}

	pat, err := m.songData.GetPatternByOrder(index.Order(order))
	if err != nil {
		if errors.Is(err, index.ErrNextPattern) {
			order++
			orderUpdated = true
			orderScanIter++
			// don't update row here
			goto orderScan
		}
		var (
			emptyRow   optional.Value[index.Row]
			emptyOrder optional.Value[index.Order]
		)
		return emptyRow, emptyOrder, err
	}

	if row >= pat.NumRows() {
		order++
		orderUpdated = true
		orderScanIter++
		row = 0
		rowUpdated = true
		goto orderScan
	}

	if orderUpdated && (forceLoopDetect || order != int(t.current.order)) && t.settings.SongLoopCount >= 0 {
		if _, found := t.songLoop.detect[index.Order(order)]; found {
			t.songLoop.current++
			if t.settings.SongLoopCount >= 0 && t.songLoop.current >= t.settings.SongLoopCount {
				var (
					emptyRow   optional.Value[index.Row]
					emptyOrder optional.Value[index.Order]
				)
				return emptyRow, emptyOrder, song.ErrStopSong
			}

			// allow and clear
			t.songLoop.detect = nil
		}
		if t.songLoop.detect == nil {
			t.songLoop.detect = make(map[index.Order]struct{})
		}
		t.songLoop.detect[index.Order(order)] = struct{}{}
	}

	var outRow optional.Value[index.Row]
	if rowUpdated {
		outRow.Set(index.Row(row))
	}

	var outOrder optional.Value[index.Order]
	if orderUpdated {
		outOrder.Set(index.Order(order))
	}

	return outRow, outOrder, nil
}
