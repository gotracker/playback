package mixing

import (
	"math"

	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/volume"
)

// PanMixerStereo is a mixer that's specialized for mixing stereo audio content
var PanMixerStereo PanMixer = &panMixerStereo{}

type panMixerStereo struct{}

func (p panMixerStereo) GetMixingMatrix(pan panning.Position, stereoSeparation float32) panning.PanMixer {
	s, c := math.Sincos(float64(pan.Angle) * 2.0)

	var d volume.Volume
	if pan.Distance > 0 {
		d = 1 / volume.Volume(pan.Distance*pan.Distance)
	}
	l := d * volume.StereoCoeff * volume.Volume(c-s)
	r := d * volume.StereoCoeff * volume.Volume(c+s)

	midCoeff := volume.Volume(1.0 / (1.0 + stereoSeparation))
	sideCoeff := volume.Volume(stereoSeparation * 0.5)

	return panning.MixerStereo{
		Matrix: volume.Matrix{
			StaticMatrix: volume.StaticMatrix{l, r},
			Channels:     2,
		},
		StereoSeparationFunc: func(in volume.Matrix) volume.Matrix {
			if stereoSeparation >= 1 {
				return in.AsStereo()
			}

			if stereoSeparation <= 0 {
				return in.AsMono().AsStereo()
			}

			dry := in.AsStereo()

			mid := (dry.StaticMatrix[0] + dry.StaticMatrix[1]) * midCoeff
			side := (dry.StaticMatrix[1] - dry.StaticMatrix[0]) * sideCoeff

			wet := dry
			wet.StaticMatrix[0] = mid - side
			wet.StaticMatrix[1] = mid + side
			return wet
		},
	}
}

func (p panMixerStereo) NumChannels() int {
	return 2
}
