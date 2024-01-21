package machine

import "fmt"

type singleRow struct {
	extraTicks int
	repeats    int
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) singleRowRowStart() error {
	var reset singleRow
	traceValueChangeWithComment(m, "extraTicks", m.extraTicks, reset.extraTicks, "RowStart")
	m.extraTicks = reset.extraTicks
	traceValueChangeWithComment(m, "repeats", m.repeats, reset.repeats, "RowStart")
	m.repeats = reset.repeats
	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) consumeRepeat() bool {
	if m.repeats <= 0 {
		return false
	}

	r := m.repeats - 1
	traceValueChangeWithComment(m, "repeats", m.repeats, r, "consumeRepeat")
	m.repeats = r
	return true
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) AddExtraTicks(ticks int) error {
	if ticks < 0 {
		return fmt.Errorf("invalid number of ticks to add: %d", ticks)
	}

	t := m.extraTicks + ticks
	traceValueChangeWithComment(m, "extraTicks", m.extraTicks, t, "AddExtraTicks")
	m.extraTicks = t
	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) RowRepeat(times int) error {
	if times < 0 {
		return fmt.Errorf("invalid number of repeat times: %d", times)
	}

	traceValueChangeWithComment(m, "repeats", m.repeats, times, "RowRepeat")
	m.repeats = times
	return nil
}
