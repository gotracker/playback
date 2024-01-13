package voice

import (
	itVolume "github.com/gotracker/playback/format/it/volume"
)

// == VolumeEnveloper ==

func (v *itVoice[TPeriod]) EnableVolumeEnvelope(enabled bool) {
	v.volEnv.SetEnabled(enabled)
}

func (v itVoice[TPeriod]) IsVolumeEnvelopeEnabled() bool {
	return v.volEnv.IsEnabled()
}

func (v itVoice[TPeriod]) GetCurrentVolumeEnvelope() itVolume.Volume {
	if v.volEnv.IsEnabled() {
		return v.volEnv.GetCurrentValue()
	}
	return itVolume.Volume(itVolume.MaxItVolume)
}

func (v *itVoice[TPeriod]) SetVolumeEnvelopePosition(pos int) {
	if doneCB := v.volEnv.SetEnvelopePosition(pos); doneCB != nil {
		doneCB(v)
	}
}

func (v itVoice[TPeriod]) GetVolumeEnvelopePosition() int {
	return v.volEnv.GetEnvelopePosition()
}
