package period

import "github.com/gotracker/playback/system"

// AmigaConverter defines a sampler period that follows the AmigaConverter-style approach of note
// definition. Useful in calculating resampling.
type AmigaConverter struct {
	System system.System
}

var _ PeriodConverter[Amiga] = (*AmigaConverter)(nil)

// GetFrequency returns the frequency defined by the period
func (c AmigaConverter) GetFrequency(p Amiga) Frequency {
	if p.IsInvalid() {
		return 0
	}
	return c.System.GetBaseClock() / Frequency(p)
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (AmigaConverter) GetSamplerAdd(p Amiga, samplerSpeed float64) float64 {
	if p.IsInvalid() {
		return 0
	}
	return samplerSpeed / float64(p)
}
