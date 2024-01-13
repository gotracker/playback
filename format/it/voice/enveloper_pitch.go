package voice

import (
	"github.com/gotracker/playback/period"
)

// == PitchEnveloper ==

func (v *itVoice[TPeriod]) EnablePitchEnvelope(enabled bool) {
	v.pitchEnv.SetEnabled(enabled)
}

func (v itVoice[TPeriod]) IsPitchEnvelopeEnabled() bool {
	if v.pitchAndFilterEnvShared && v.filterEnvActive {
		return false
	}
	return v.pitchEnv.IsEnabled()
}

func (v itVoice[TPeriod]) GetCurrentPitchEnvelope() period.Delta {
	if v.pitchEnv.IsEnabled() {
		return v.pitchEnv.GetCurrentValue()
	}
	return 0
}

func (v *itVoice[TPeriod]) SetPitchEnvelopePosition(pos int) {
	if !v.pitchAndFilterEnvShared || !v.filterEnvActive {
		if doneCB := v.pitchEnv.SetEnvelopePosition(pos); doneCB != nil {
			doneCB(v)
		}
	}
}

func (v itVoice[TPeriod]) GetPitchEnvelopePosition() int {
	return v.pitchEnv.GetEnvelopePosition()
}
