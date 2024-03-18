package voice

import (
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/voice/types"
)

// == AmpModulator ==

func (v *xmVoice[TPeriod]) SetActive(on bool) error {
	return v.amp.SetActive(on)
}

func (v xmVoice[TPeriod]) IsActive() bool {
	return v.amp.IsActive()
}

func (v *xmVoice[TPeriod]) SetMixingVolume(vol xmVolume.XmVolume) error {
	return v.amp.SetMixingVolume(vol)
}

func (v xmVoice[TPeriod]) GetMixingVolume() xmVolume.XmVolume {
	return v.amp.GetMixingVolume()
}

func (v *xmVoice[TPeriod]) SetVolume(vol xmVolume.XmVolume) error {
	if vol.IsUseInstrumentVol() {
		vol = v.voicer.GetDefaultVolume()
	}
	return v.amp.SetVolume(vol)
}

func (v xmVoice[TPeriod]) GetVolume() xmVolume.XmVolume {
	return v.amp.GetVolume()
}

func (v *xmVoice[TPeriod]) SetVolumeDelta(d types.VolumeDelta) error {
	return v.amp.SetVolumeDelta(d)
}

func (v xmVoice[TPeriod]) GetVolumeDelta() types.VolumeDelta {
	return v.amp.GetVolumeDelta()
}

func (v xmVoice[TPeriod]) IsFadeout() bool {
	return v.fadeout.IsActive()
}

func (v xmVoice[TPeriod]) GetFadeoutVolume() volume.Volume {
	return v.fadeout.GetVolume()
}

func (v xmVoice[TPeriod]) GetFinalVolume() volume.Volume {
	vol := v.amp.GetFinalVolume()
	if v.IsVolumeEnvelopeEnabled() {
		vol *= v.GetCurrentVolumeEnvelope().ToVolume()
	}
	return vol * v.fadeout.GetFinalVolume()
}
