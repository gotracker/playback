package sampler

import (
	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/playback/period"
)

// Sampler is a container of sampler/mixer settings
type Sampler struct {
	SampleRate    int
	BaseClockRate period.Frequency

	mixer mixing.Mixer
}

// NewSampler returns a new sampler object based on the input settings
func NewSampler(samplesPerSec, channels int) *Sampler {
	s := Sampler{
		SampleRate: samplesPerSec,
		mixer: mixing.Mixer{
			Channels: channels,
		},
	}
	return &s
}

// Mixer returns a pointer to the current mixer object
func (s *Sampler) Mixer() *mixing.Mixer {
	return &s.mixer
}

// GetPanMixer returns the panning mixer that can generate a matrix
// based on input pan value
func (s *Sampler) GetPanMixer() mixing.PanMixer {
	return mixing.GetPanMixer(s.mixer.Channels)
}
