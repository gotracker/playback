package component

import (
	"fmt"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/tracing"
)

// FreqModulator is a frequency (pitch) modulator
type FreqModulator[TPeriod period.Period] struct {
	settings FreqModulatorSettings[TPeriod]
	period   TPeriod
	delta    period.Delta
}

type FreqModulatorSettings[TPeriod period.Period] struct {
}

func (f *FreqModulator[TPeriod]) Setup(settings FreqModulatorSettings[TPeriod]) {
	f.settings = settings
}

func (f FreqModulator[TPeriod]) Clone() FreqModulator[TPeriod] {
	m := f
	return m
}

// SetPeriod sets the current period (before AutoVibrato and Delta calculation)
func (f *FreqModulator[TPeriod]) SetPeriod(period TPeriod) {
	f.period = period
}

// GetPeriod returns the current period (before AutoVibrato and Delta calculation)
func (f *FreqModulator[TPeriod]) GetPeriod() TPeriod {
	return f.period
}

// SetPeriodDelta sets the current period delta (before AutoVibrato calculation)
func (f *FreqModulator[TPeriod]) SetPeriodDelta(delta period.Delta) {
	f.delta = delta
}

// GetDelta returns the current period delta (before AutoVibrato calculation)
func (f *FreqModulator[TPeriod]) GetPeriodDelta() period.Delta {
	return f.delta
}

// GetFinalPeriod returns the current period (after AutoVibrato and Delta calculation)
func (f *FreqModulator[TPeriod]) GetFinalPeriod() TPeriod {
	return period.AddDelta(f.period, f.delta)
}

func (f FreqModulator[TPeriod]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("period{%v} delta{%v}",
		f.period,
		f.delta,
	), comment)
}
