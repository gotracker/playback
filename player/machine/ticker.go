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

	row, order, err := advanceRowOrder(t, m)
	if err != nil {
		return err
	}

	t.current.row = row
	t.current.order = order
	m.us.SetTracingTick(t.current.order, t.current.row, t.current.tick)

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
	done := false
	if tick >= m.tempo {
		tick = 0
		rowAdvanced = true
	}

	if rowAdvanced {
		if err := m.onRowEnd(); err != nil {
			return err
		}

		var err error
		row, order, err = advanceRowOrder(t, m)
		if err != nil && !errors.Is(err, song.ErrStopSong) {
			return err
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
	} else if rowAdvanced {
		if err := m.onRowStart(); err != nil {
			return err
		}
	}
	return nil
}

func advanceRowOrder[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](t *ticker, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) (index.Row, index.Order, error) {
	row := int(t.current.row)
	order := int(t.current.order)

	desiredRow, rowSet := t.next.row.Get()
	desiredOrder, orderSet := t.next.order.Get()

	if rowSet && orderSet {
		row = int(desiredRow)
		order = int(desiredOrder)
	} else if rowSet {
		row = int(desiredRow)
		order++
	} else if orderSet {
		order = int(desiredOrder)
		row = 0
	} else {
		row++
	}

	t.next.row.Reset()
	t.next.order.Reset()

	orderScanMax := len(m.songData.GetOrderList())
	orderScanIter := 0
orderScan:
	if orderScanIter >= orderScanMax {
		return index.Row(row), index.Order(order), song.ErrStopSong
	}

	pat, err := m.songData.GetPatternIntfByOrder(index.Order(order))
	if err != nil {
		if errors.Is(err, index.ErrNextPattern) {
			order++
			orderScanIter++
			// don't update row here
			goto orderScan
		}
		return index.Row(row), index.Order(order), err
	}

	if row >= pat.NumRows() {
		order++
		orderScanIter++
		row = 0
		goto orderScan
	}

	return index.Row(row), index.Order(order), nil
}
