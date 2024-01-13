package voice

import (
	xmPanning "github.com/gotracker/playback/format/xm/panning"
)

// == PanEnveloper ==

func (v *xmVoice[TPeriod]) EnablePanEnvelope(enabled bool) {
	v.panEnv.SetEnabled(enabled)
}

func (v xmVoice[TPeriod]) IsPanEnvelopeEnabled() bool {
	return v.panEnv.IsEnabled()
}

func (v xmVoice[TPeriod]) GetCurrentPanEnvelope() xmPanning.Panning {
	if v.panEnv.IsEnabled() {
		return v.panEnv.GetCurrentValue()
	}
	return xmPanning.DefaultPanning
}

func (v *xmVoice[TPeriod]) SetPanEnvelopePosition(pos int) {
	if doneCB := v.panEnv.SetEnvelopePosition(pos); doneCB != nil {
		doneCB(v)
	}
}

func (v xmVoice[TPeriod]) GetPanEnvelopePosition() int {
	return v.panEnv.GetEnvelopePosition()
}
