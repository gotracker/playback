package period

import (
	"github.com/gotracker/playback/period"
)

// AmigaConverter defines a sampler period that follows the AmigaConverter-style approach of note
// definition. Useful in calculating resampling.
type amigaConverter struct{}

var AmigaConverter period.PeriodConverter[period.Amiga] = amigaConverter{}

// GetFrequency returns the frequency defined by the period
func (amigaConverter) GetFrequency(p period.Amiga) period.Frequency {
	if p.IsInvalid() {
		return 0
	}
	return period.Frequency(ITBaseClock) / period.Frequency(p)
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (amigaConverter) GetSamplerAdd(p period.Amiga, samplerSpeed float64) float64 {
	if p.IsInvalid() {
		return 0
	}
	return samplerSpeed / float64(p)
}
