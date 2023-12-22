package period

import (
	"math"

	"github.com/gotracker/playback/note"
)

// LinearConverter defines a sampler period that follows the Linear-style approach of note
// definition. Useful in calculating resampling.
type LinearConverter struct {
	BaseClock      Frequency
	BaseFinetune   note.Finetune
	FinesPerOctave note.Finetune
}

var _ PeriodConverter[Linear] = (*LinearConverter)(nil)

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (c LinearConverter) GetSamplerAdd(p Linear, samplerSpeed float64) float64 {
	return float64(c.GetFrequency(p)) * samplerSpeed / float64(c.BaseClock)
}

// GetFrequency returns the frequency defined by the period
func (c LinearConverter) GetFrequency(p Linear) Frequency {
	pft := float64(p.Finetune-c.BaseFinetune) / float64(c.FinesPerOctave)
	f := p.CommonRate * Frequency(math.Pow(2.0, pft))
	return f
}
