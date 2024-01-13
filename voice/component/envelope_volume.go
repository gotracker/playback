package component

import (
	"github.com/gotracker/playback/util"
	"github.com/gotracker/playback/voice/types"
)

// VolumeEnvelope is an amplitude modulation envelope
type VolumeEnvelope[TVolume types.Volume] struct {
	baseEnvelope[TVolume, TVolume]
}

func (e *VolumeEnvelope[TVolume]) Setup(settings EnvelopeSettings[TVolume, TVolume]) {
	e.baseEnvelope.Setup(settings, e.calc)
}

func (e VolumeEnvelope[TVolume]) Clone() VolumeEnvelope[TVolume] {
	var m VolumeEnvelope[TVolume]
	m.baseEnvelope = e.baseEnvelope.Clone(m.calc)
	return m
}

func (e *VolumeEnvelope[TVolume]) calc() TVolume {
	cur, next, t := e.state.GetCurrentValue(e.keyOn, e.prevKeyOn)

	var y0 TVolume
	if cur != nil {
		y0 = cur.Y
	}

	var y1 TVolume
	if next != nil {
		y1 = next.Y
	}

	return util.Lerp(float64(t), y0, y1)
}
