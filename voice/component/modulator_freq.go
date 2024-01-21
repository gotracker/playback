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
	PC period.PeriodConverter[TPeriod]
}

func (f *FreqModulator[TPeriod]) Setup(settings FreqModulatorSettings[TPeriod]) error {
	f.settings = settings
	var empty TPeriod
	f.unkeyed.period = empty
	return f.Reset()
}

func (f *FreqModulator[TPeriod]) Reset() error {
	f.keyed.delta = 0
	return f.updateFinal()
}

func (f FreqModulator[TPeriod]) Clone() FreqModulator[TPeriod] {
	m := f
	return m
}

// SetPeriod sets the current period (before AutoVibrato and Delta calculation)
func (f *FreqModulator[TPeriod]) SetPeriod(period TPeriod) error {
	if period.IsInvalid() {
		// ignore it for now
		return nil
	}

	f.unkeyed.period = period
	return f.updateFinal()
}

// GetPeriod returns the current period (before AutoVibrato and Delta calculation)
func (f *FreqModulator[TPeriod]) GetPeriod() TPeriod {
	return f.unkeyed.period
}

// SetPeriodDelta sets the current period delta (before AutoVibrato calculation)
func (f *FreqModulator[TPeriod]) SetPeriodDelta(delta period.Delta) error {
	f.keyed.delta = delta
	return f.updateFinal()
}

// GetDelta returns the current period delta (before AutoVibrato calculation)
func (f *FreqModulator[TPeriod]) GetPeriodDelta() period.Delta {
	return f.keyed.delta
}

// GetFinalPeriod returns the current period (after AutoVibrato and Delta calculation)
func (f *FreqModulator[TPeriod]) GetFinalPeriod() (TPeriod, error) {
	return f.final, nil
}

func (f FreqModulator[TPeriod]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("period{%v} delta{%v} final{%v}",
		f.unkeyed.period,
		f.keyed.delta,
		f.final,
	), comment)
}

func (f *FreqModulator[TPeriod]) updateFinal() error {
	var err error
	f.final, err = f.settings.PC.AddDelta(f.unkeyed.period, f.keyed.delta)
	return err
}
