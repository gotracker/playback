package system

import "github.com/gotracker/playback/frequency"

type System interface {
	GetMaxPastNotesPerChannel() int
	GetCommonRate() frequency.Frequency
	GetSamplerSpeed(sampleRate frequency.Frequency) float32
}
