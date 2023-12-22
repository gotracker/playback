package component

import (
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/oscillator"
)

// FreqModulator is a frequency (pitch) modulator
type FreqModulator[TPeriod period.Period] struct {
	period             TPeriod
	delta              period.Delta
	autoVibratoEnabled bool
	autoVibrato        oscillator.Oscillator
	autoVibratoDepth   float32
	autoVibratoRate    int
	autoVibratoSweep   int // maximum age when oscillator is at max depth (in ticks)
	autoVibratoAge     int // current age of oscillator (in ticks)
}

// SetPeriod sets the current period (before AutoVibrato and Delta calculation)
func (a *FreqModulator[TPeriod]) SetPeriod(period TPeriod) {
	a.period = period
}

// GetPeriod returns the current period (before AutoVibrato and Delta calculation)
func (a *FreqModulator[TPeriod]) GetPeriod() TPeriod {
	return a.period
}

// SetDelta sets the current period delta (before AutoVibrato calculation)
func (a *FreqModulator[TPeriod]) SetDelta(delta period.Delta) {
	a.delta = delta
}

// GetDelta returns the current period delta (before AutoVibrato calculation)
func (a *FreqModulator[TPeriod]) GetDelta() period.Delta {
	return a.delta
}

// SetAutoVibratoEnabled sets the status of the AutoVibrato enablement flag
func (a *FreqModulator[TPeriod]) SetAutoVibratoEnabled(enabled bool) {
	a.autoVibratoEnabled = enabled
}

// ConfigureAutoVibrato sets the AutoVibrato oscillator settings
func (a *FreqModulator[TPeriod]) ConfigureAutoVibrato(av voice.AutoVibrato) {
	a.autoVibrato = av.Generate()
	a.autoVibratoRate = int(av.Rate)
	a.autoVibratoDepth = av.Depth
}

// ResetAutoVibrato resets the current AutoVibrato
func (a *FreqModulator[TPeriod]) ResetAutoVibrato(sweep ...int) {
	if a.autoVibrato != nil {
		a.autoVibrato.Reset(true)
	}

	a.autoVibratoAge = 0

	if sweep != nil {
		a.autoVibratoSweep = sweep[0]
	}
}

// IsAutoVibratoEnabled returns the status of the AutoVibrato enablement flag
func (a *FreqModulator[TPeriod]) IsAutoVibratoEnabled() bool {
	return a.autoVibratoEnabled
}

// GetFinalPeriod returns the current period (after AutoVibrato and Delta calculation)
func (a *FreqModulator[TPeriod]) GetFinalPeriod() TPeriod {
	p := period.AddDelta(a.period, a.delta)
	if a.autoVibratoEnabled {
		depth := a.autoVibratoDepth
		if a.autoVibratoSweep > a.autoVibratoAge {
			depth *= float32(a.autoVibratoAge) / float32(a.autoVibratoSweep)
		}
		avDelta := a.autoVibrato.GetWave(depth)
		d := period.Delta(avDelta)
		p = period.AddDelta(p, d)
	}

	return p
}

// Advance advances the autoVibrato value by 1 tick
func (a *FreqModulator[TPeriod]) Advance() {
	if !a.autoVibratoEnabled {
		return
	}

	a.autoVibrato.Advance(a.autoVibratoRate)
	a.autoVibratoAge++
}
