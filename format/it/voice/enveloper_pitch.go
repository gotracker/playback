package voice

import (
	"github.com/gotracker/playback/period"
)

// == PitchEnveloper ==

func (v *itVoice[TPeriod]) EnablePitchEnvelope(enabled bool) error {
	return v.pitchEnv.SetEnabled(enabled)
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

func (v *itVoice[TPeriod]) SetPitchEnvelopePosition(pos int) error {
	if !v.pitchAndFilterEnvShared || !v.filterEnvActive {
		doneCB, err := v.pitchEnv.SetEnvelopePosition(pos)
		if err != nil {
			return err
		}
		if doneCB != nil {
			doneCB(v)
		}
	}
	return nil
}

func (v itVoice[TPeriod]) GetPitchEnvelopePosition() int {
	return v.pitchEnv.GetEnvelopePosition()
}
