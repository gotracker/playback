package component

import (
	"github.com/gotracker/playback/util"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/types"
)

// VolumeEnvelope is an amplitude modulation envelope
type VolumeEnvelope[TVolume types.Volume] struct {
	baseEnvelope[TVolume, TVolume]
}

func (e *VolumeEnvelope[TVolume]) Setup(settings EnvelopeSettings[TVolume, TVolume]) {
	e.baseEnvelope.Setup(settings, e.calc)
}

func (e VolumeEnvelope[TVolume]) Clone(onFinished voice.Callback) VolumeEnvelope[TVolume] {
	var m VolumeEnvelope[TVolume]
	m.baseEnvelope = e.baseEnvelope.Clone(m.calc, onFinished)
	return m
}

func (e *VolumeEnvelope[TVolume]) calc(y0, y1 TVolume, t float64) TVolume {
	return util.Lerp(t, y0, y1)
}
