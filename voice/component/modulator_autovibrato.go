package component

import (
	"fmt"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/autovibrato"
	"github.com/gotracker/playback/voice/oscillator"
)

// AutoVibratoModulator is a frequency (pitch) modulator
type AutoVibratoModulator[TPeriod period.Period] struct {
	settings    autovibrato.AutoVibratoSettings
	enabled     bool
	autoVibrato oscillator.Oscillator
	age         int // current age of oscillator (in ticks)
}

func (f *AutoVibratoModulator[TPeriod]) Setup(settings autovibrato.AutoVibratoSettings) {
	f.settings = settings
	f.Reset()
}

func (f AutoVibratoModulator[TPeriod]) Clone() AutoVibratoModulator[TPeriod] {
	m := f
	if f.autoVibrato != nil {
		m.autoVibrato = f.autoVibrato.Clone()
	}
	return m
}

func (f *AutoVibratoModulator[TPeriod]) Reset() {
	f.enabled = f.settings.Enabled
	f.autoVibrato = f.settings.Generate()
	f.ResetAutoVibrato()
}

// SetEnabled sets the status of the AutoVibrato enablement flag
func (f *AutoVibratoModulator[TPeriod]) SetEnabled(enabled bool) {
	f.enabled = enabled
}

// ConfigureAutoVibrato sets the AutoVibrato oscillator settings
func (f *AutoVibratoModulator[TPeriod]) ConfigureAutoVibrato() {
	f.autoVibrato = f.settings.Generate()
}

// ResetAutoVibrato resets the current AutoVibrato
func (f *AutoVibratoModulator[TPeriod]) ResetAutoVibrato() {
	if f.autoVibrato != nil {
		f.autoVibrato.HardReset()
	}

	f.age = 0
}

// IsAutoVibratoEnabled returns the status of the AutoVibrato enablement flag
func (f *AutoVibratoModulator[TPeriod]) IsAutoVibratoEnabled() bool {
	return f.enabled
}

// GetFinalPeriod returns the current period (after AutoVibrato and Delta calculation)
func (f *AutoVibratoModulator[TPeriod]) GetAdjustedPeriod(in TPeriod) TPeriod {
	if !f.enabled {
		return in
	}

	depth := f.settings.Depth
	if f.settings.Sweep > f.age {
		depth *= float32(f.age) / float32(f.settings.Sweep)
	}
	avDelta := f.autoVibrato.GetWave(depth)
	d := period.Delta(avDelta)
	return period.AddDelta(in, d)
}

// Advance advances the autoVibrato value by 1 tick
func (f *AutoVibratoModulator[TPeriod]) Advance() {
	if !f.enabled {
		return
	}

	f.autoVibrato.Advance(f.settings.Rate)
	f.age++
}

func (f AutoVibratoModulator[TPeriod]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("enabled{%v} age{%v}",
		f.enabled,
		f.age,
	), comment)
}
