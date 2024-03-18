package mixing

import (
	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/volume"
)

// PanMixerMono is a mixer that's specialized for mixing monaural audio content
var PanMixerMono PanMixer = &panMixerMono{}

type panMixerMono struct{}

func (p panMixerMono) GetMixingMatrix(pan panning.Position, stereoSeparation float32) panning.PanMixer {
	// distance and angle are ignored on mono
	return volume.Matrix{
		StaticMatrix: volume.StaticMatrix{1.0},
		Channels:     1,
	}
}

func (p panMixerMono) NumChannels() int {
	return 1
}
