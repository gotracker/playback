package system

type System interface {
	GetMaxPastNotesPerChannel() int
	GetCommonRate() Frequency
	GetSamplerSpeed(sampleRate Frequency) float32
}
