package voice

import (
	itPanning "github.com/gotracker/playback/format/it/panning"
)

// == PanEnveloper ==

func (v *itVoice[TPeriod]) EnablePanEnvelope(enabled bool) {
	v.panEnv.SetEnabled(enabled)
}

func (v itVoice[TPeriod]) IsPanEnvelopeEnabled() bool {
	return v.panEnv.IsEnabled()
}

func (v itVoice[TPeriod]) GetCurrentPanEnvelope() itPanning.Panning {
	if v.panEnv.IsEnabled() {
		return v.panEnv.GetCurrentValue()
	}
	return itPanning.DefaultPanning
}

func (v *itVoice[TPeriod]) SetPanEnvelopePosition(pos int) {
	if doneCB := v.panEnv.SetEnvelopePosition(pos); doneCB != nil {
		doneCB(v)
	}
}

func (v itVoice[TPeriod]) GetPanEnvelopePosition() int {
	return v.panEnv.GetEnvelopePosition()
}
