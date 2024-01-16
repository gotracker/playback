package voice

import (
	xmVolume "github.com/gotracker/playback/format/xm/volume"
)

// == VolumeEnveloper ==

func (v *xmVoice[TPeriod]) EnableVolumeEnvelope(enabled bool) error {
	return v.volEnv.SetEnabled(enabled)
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

func (v *xmVoice[TPeriod]) SetVolumeEnvelopePosition(pos int) error {
	doneCB, err := v.volEnv.SetEnvelopePosition(pos)
	if err != nil {
		return err
	}
	if doneCB != nil {
		doneCB(v)
	}
	return nil
}

func (v xmVoice[TPeriod]) GetVolumeEnvelopePosition() int {
	return v.volEnv.GetEnvelopePosition()
}
