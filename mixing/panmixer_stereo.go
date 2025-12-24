package mixing

import (
	"math"

	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/volume"
)

// PanMixerStereo is a mixer that's specialized for mixing stereo audio content
var PanMixerStereo PanMixer = &panMixerStereo{}

type panMixerStereo struct{}

// Dolby Pro Logic II surround encoding: feed surrounds as equal-magnitude,
// opposite-phase signals to L/R (nominal +90/-90deg). We approximate using
// out-of-phase amplitudes and bypass additional stereo separation to
// preserve phase cues for downstream decoders.
const surroundScale = volume.Volume(1 / math.Sqrt2)

func (p panMixerStereo) GetMixingMatrix(pan panning.Position, stereoSeparation float32) panning.PanMixer {
	if pan == panning.SurroundPosition {
		return panning.MixerStereo{
			Matrix: volume.Matrix{
				StaticMatrix: volume.StaticMatrix{surroundScale, -surroundScale},
				Channels:     2,
			},
			StereoSeparationFunc: func(in volume.Matrix) volume.Matrix {
				return in.AsStereo()
			},
		}
	}

	s, c := math.Sincos(float64(pan.Angle))

	var d volume.Volume
	if pan.Distance > 0 {
		d = 1 / volume.Volume(pan.Distance*pan.Distance)
	}
	l := d * volume.Volume(s)
	r := d * volume.Volume(c)

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
