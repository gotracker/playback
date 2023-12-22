package period

import (
	"math"

	"github.com/gotracker/playback/period"
)

// Linear is a linear period, based on semitone and finetune values
type linearConverter struct{}

var LinearConverter period.PeriodConverter[period.Linear] = linearConverter{}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (c linearConverter) GetSamplerAdd(p period.Linear, samplerSpeed float64) float64 {
	return float64(c.GetFrequency(p)) * samplerSpeed / float64(ITBaseClock)
}

// GetFrequency returns the frequency defined by the period
func (linearConverter) GetFrequency(p period.Linear) period.Frequency {
	pft := float64(p.Finetune-C5SlideFines) / float64(SlideFinesPerOctave)
	f := p.CommonRate * period.Frequency(math.Pow(2.0, pft))
	return f
}
