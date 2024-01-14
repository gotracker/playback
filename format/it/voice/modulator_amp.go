package voice

import (
	"github.com/gotracker/gomixing/volume"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/voice/types"
)

// == AmpModulator ==

func (v *itVoice[TPeriod]) SetActive(on bool) {
	v.amp.SetActive(on)
}

func (v itVoice[TPeriod]) IsActive() bool {
	return v.amp.IsActive()
}

func (v *itVoice[TPeriod]) SetMixingVolume(vol itVolume.FineVolume) {
	if !vol.IsUseInstrumentVol() {
		v.amp.SetMixingVolume(vol)
	}
}

func (v itVoice[TPeriod]) GetMixingVolume() itVolume.FineVolume {
	return v.amp.GetMixingVolume()
}

func (v *itVoice[TPeriod]) SetVolume(vol itVolume.Volume) {
	if vol.IsUseInstrumentVol() {
		vol = v.voicer.GetDefaultVolume()
	}
	v.amp.SetVolume(vol)
}

func (v itVoice[TPeriod]) GetVolume() itVolume.Volume {
	return v.amp.GetVolume()
}

func (v *itVoice[TPeriod]) SetVolumeDelta(d types.VolumeDelta) {
	v.amp.SetVolumeDelta(d)
}

func (v itVoice[TPeriod]) GetVolumeDelta() types.VolumeDelta {
	return v.amp.GetVolumeDelta()
}

func (v itVoice[TPeriod]) IsFadeout() bool {
	return v.fadeout.IsActive()
}

func (v itVoice[TPeriod]) GetFadeoutVolume() volume.Volume {
	return v.fadeout.GetVolume()
}

func (v itVoice[TPeriod]) GetFinalVolume() volume.Volume {
	return v.finalVol
}
