package mixing

import (
	"math"

	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/volume"
)

// PanMixerQuad is a mixer that's specialized for mixing quadraphonic audio content
var PanMixerQuad PanMixer = &panMixerQuad{}

type panMixerQuad struct{}

func (p panMixerQuad) GetMixingMatrix(pan panning.Position, stereoSeparation float32) panning.PanMixer {
	pangle := float64(pan.Angle)
	sf, cf := math.Sincos(pangle)
	sr, cr := math.Sin(pangle+math.Pi/2.0), math.Cos(pangle-math.Pi/2.0)
	var d volume.Volume
	if pan.Distance > 0 {
		d = 1 / volume.Volume(pan.Distance*pan.Distance)
	}
	lf := d * volume.Volume(sf)
	rf := d * volume.Volume(cf)
	lr := d * volume.Volume(cr)
	rr := d * volume.Volume(sr)
	return volume.Matrix{
		StaticMatrix: volume.StaticMatrix{lf, rf, lr, rr},
		Channels:     4,
	}
}

func (p panMixerQuad) NumChannels() int {
	return 4
}
