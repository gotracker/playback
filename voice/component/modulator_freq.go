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
	unkeyed  struct {
		period TPeriod
	}
	keyed struct {
		delta period.Delta
	}
	final TPeriod
}

type FreqModulatorSettings[TPeriod period.Period] struct {
}

func (f *FreqModulator[TPeriod]) Setup(settings FreqModulatorSettings[TPeriod]) {
	f.settings = settings
	var empty TPeriod
	f.unkeyed.period = empty
	f.Reset()
}

func (f *FreqModulator[TPeriod]) Reset() {
	f.keyed.delta = 0
	f.updateFinal()
}

func (f FreqModulator[TPeriod]) Clone() FreqModulator[TPeriod] {
	m := f
	return m
}

// SetPeriod sets the current period (before AutoVibrato and Delta calculation)
func (f *FreqModulator[TPeriod]) SetPeriod(period TPeriod) {
	if period.IsInvalid() {
		return
	}

	f.unkeyed.period = period
	f.updateFinal()
}

// GetPeriod returns the current period (before AutoVibrato and Delta calculation)
func (f *FreqModulator[TPeriod]) GetPeriod() TPeriod {
	return f.unkeyed.period
}

// SetPeriodDelta sets the current period delta (before AutoVibrato calculation)
func (f *FreqModulator[TPeriod]) SetPeriodDelta(delta period.Delta) {
	f.keyed.delta = delta
	f.updateFinal()
}

// GetDelta returns the current period delta (before AutoVibrato calculation)
func (f *FreqModulator[TPeriod]) GetPeriodDelta() period.Delta {
	return f.keyed.delta
}

// GetFinalPeriod returns the current period (after AutoVibrato and Delta calculation)
func (f *FreqModulator[TPeriod]) GetFinalPeriod() TPeriod {
	return f.final
}

func (f FreqModulator[TPeriod]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("period{%v} delta{%v} final{%v}",
		f.unkeyed.period,
		f.keyed.delta,
		f.final,
	), comment)
}

func (f *FreqModulator[TPeriod]) updateFinal() {
	f.final = period.AddDelta(f.unkeyed.period, f.keyed.delta)
}
