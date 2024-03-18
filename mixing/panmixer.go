package mixing

import (
	"github.com/gotracker/playback/mixing/panning"
)

// PanMixer is a mixer that's specialized for mixing multichannel audio content
type PanMixer interface {
	GetMixingMatrix(pan panning.Position, stereoSeparation float32) panning.PanMixer
	NumChannels() int
}

// GetPanMixer returns the panning mixer that can generate a matrix
// based on input pan value
func GetPanMixer(channels int) PanMixer {
	switch channels {
	case 1:
		return PanMixerMono
	case 2:
		return PanMixerStereo
	case 4:
		return PanMixerQuad
	}

	return nil
}
