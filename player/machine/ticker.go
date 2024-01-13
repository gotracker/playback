package machine

import (
	"errors"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/song"
	"github.com/heucuva/optional"
)

type ticker struct {
	current struct {
		tick  int
		row   index.Row
		order index.Order
	}
	next struct {
		row   optional.Value[index.Row]
		order optional.Value[index.Order]
	}
}

type tickerSettings struct {
	Order index.Order
	Row   index.Row
}

func initTick[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](t *ticker, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], setup tickerSettings) error {
	t.current.tick = 0
	t.current.row = 0
	t.current.order = 0
	t.next.row.Set(setup.Row)
	t.next.order.Set(setup.Order)

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
		row = int(desiredRow)
		rowUpdated = true
		order = int(desiredOrder)
		orderUpdated = true
	} else if rowSet {
		row = int(desiredRow)
		rowUpdated = true
		order++
		orderUpdated = true
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
orderScan:
	if orderScanIter >= orderScanMax {
		var (
			emptyRow   optional.Value[index.Row]
			emptyOrder optional.Value[index.Order]
		)
		return emptyRow, emptyOrder, song.ErrStopSong
	}

	pat, err := m.songData.GetPatternIntfByOrder(index.Order(order))
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
