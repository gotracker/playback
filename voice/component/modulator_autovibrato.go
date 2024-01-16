package component

import (
	"fmt"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/autovibrato"
	"github.com/gotracker/playback/voice/oscillator"
	"github.com/gotracker/playback/voice/types"
)

// AutoVibratoModulator is a frequency (pitch) modulator
type AutoVibratoModulator[TPeriod types.Period] struct {
	settings autovibrato.AutoVibratoSettings[TPeriod]
	unkeyed  struct {
		enabled bool
	}
	keyed struct {
		age int // current age of oscillator (in ticks)
	}
	autoVibrato oscillator.Oscillator
}

func (f *AutoVibratoModulator[TPeriod]) Setup(settings autovibrato.AutoVibratoSettings[TPeriod]) {
	f.settings = settings
	f.unkeyed.enabled = f.settings.Enabled
	f.Reset()
}

func (f AutoVibratoModulator[TPeriod]) Clone() AutoVibratoModulator[TPeriod] {
	m := f
	if f.autoVibrato != nil {
		m.autoVibrato = f.autoVibrato.Clone()
	}
	return m
}

func (f *AutoVibratoModulator[TPeriod]) Reset() error {
	f.keyed.age = 0
	f.autoVibrato = f.settings.Generate()
	return f.ResetAutoVibrato()
}

// SetEnabled sets the status of the AutoVibrato enablement flag
func (f *AutoVibratoModulator[TPeriod]) SetEnabled(enabled bool) {
	f.unkeyed.enabled = enabled
}

// ConfigureAutoVibrato sets the AutoVibrato oscillator settings
func (f *AutoVibratoModulator[TPeriod]) ConfigureAutoVibrato() {
	f.autoVibrato = f.settings.Generate()
}

// ResetAutoVibrato resets the current AutoVibrato
func (f *AutoVibratoModulator[TPeriod]) ResetAutoVibrato() error {
	if f.autoVibrato != nil {
		f.autoVibrato.HardReset()
	}

	f.keyed.age = 0
	return nil
}

// IsAutoVibratoEnabled returns the status of the AutoVibrato enablement flag
func (f *AutoVibratoModulator[TPeriod]) IsAutoVibratoEnabled() bool {
	return f.unkeyed.enabled
}

// GetFinalPeriod returns the current period (after AutoVibrato and Delta calculation)
func (f *AutoVibratoModulator[TPeriod]) GetAdjustedPeriod(in TPeriod) (TPeriod, error) {
	if !f.unkeyed.enabled {
		return in, nil
	}

	depth := f.settings.Depth
	if f.settings.Sweep > f.keyed.age {
		depth *= float32(f.keyed.age) / float32(f.settings.Sweep)
	}
	avDelta := f.autoVibrato.GetWave(depth)
	d := period.Delta(avDelta)
	return f.settings.PC.AddDelta(in, d)
}

// Advance advances the autoVibrato value by 1 tick
func (f *AutoVibratoModulator[TPeriod]) Advance() {
	if !f.unkeyed.enabled {
		return
	}

	f.autoVibrato.Advance(f.settings.Rate)
	f.keyed.age++
}

func (f AutoVibratoModulator[TPeriod]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("enabled{%v} age{%v}",
		f.unkeyed.enabled,
		f.keyed.age,
	), comment)
}
