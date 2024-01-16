package voice

import (
	xmPanning "github.com/gotracker/playback/format/xm/panning"
)

// == PanEnveloper ==

func (v *xmVoice[TPeriod]) EnablePanEnvelope(enabled bool) error {
	return v.panEnv.SetEnabled(enabled)
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

func (v *xmVoice[TPeriod]) SetPanEnvelopePosition(pos int) error {
	doneCB, err := v.panEnv.SetEnvelopePosition(pos)
	if err != nil {
		return err
	}
	if doneCB != nil {
		doneCB(v)
	}
	return nil
}

func (v xmVoice[TPeriod]) GetPanEnvelopePosition() int {
	return v.panEnv.GetEnvelopePosition()
}
