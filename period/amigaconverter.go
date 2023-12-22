package period

// AmigaConverter defines a sampler period that follows the AmigaConverter-style approach of note
// definition. Useful in calculating resampling.
type AmigaConverter struct {
	BaseClock Frequency
}

var _ PeriodConverter[Amiga] = (*AmigaConverter)(nil)

// GetFrequency returns the frequency defined by the period
func (c AmigaConverter) GetFrequency(p Amiga) Frequency {
	if p.IsInvalid() {
		return 0
	}
	return c.BaseClock / Frequency(p)
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (AmigaConverter) GetSamplerAdd(p Amiga, samplerSpeed float64) float64 {
	if p.IsInvalid() {
		return 0
	}
	return samplerSpeed / float64(p)
}
