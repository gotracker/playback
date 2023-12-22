package period

type PeriodConverter[TPeriod Period] interface {
	GetSamplerAdd(TPeriod, float64) float64
	GetFrequency(TPeriod) Frequency
}
