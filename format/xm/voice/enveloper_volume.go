package voice

import (
	xmVolume "github.com/gotracker/playback/format/xm/volume"
)

// == VolumeEnveloper ==

func (v *xmVoice[TPeriod]) EnableVolumeEnvelope(enabled bool) {
	v.volEnv.SetEnabled(enabled)
}

func (v xmVoice[TPeriod]) IsVolumeEnvelopeEnabled() bool {
	return v.volEnv.IsEnabled()
}

func (v xmVoice[TPeriod]) GetCurrentVolumeEnvelope() xmVolume.XmVolume {
	if v.volEnv.IsEnabled() {
		return v.volEnv.GetCurrentValue()
	}
	return xmVolume.DefaultXmVolume
}

func (v *xmVoice[TPeriod]) SetVolumeEnvelopePosition(pos int) {
	if doneCB := v.volEnv.SetEnvelopePosition(pos); doneCB != nil {
		doneCB(v)
	}
}

func (v xmVoice[TPeriod]) GetVolumeEnvelopePosition() int {
	return v.volEnv.GetEnvelopePosition()
}
