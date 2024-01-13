package voice

// == FilterEnveloper ==

func (v *itVoice[TPeriod]) EnableFilterEnvelope(enabled bool) {
	if !v.pitchAndFilterEnvShared {
		v.filterEnv.SetEnabled(enabled)
		return
	}

	// shared filter/pitch envelope
	if !v.filterEnvActive {
		return
	}

	v.filterEnv.SetEnabled(enabled)
}

func (v itVoice[TPeriod]) IsFilterEnvelopeEnabled() bool {
	if v.pitchAndFilterEnvShared && !v.filterEnvActive {
		return false
	}
	return v.filterEnv.IsEnabled()
}

func (v itVoice[TPeriod]) GetCurrentFilterEnvelope() uint8 {
	return v.filterEnv.GetCurrentValue()
}

func (v *itVoice[TPeriod]) SetFilterEnvelopePosition(pos int) {
	if !v.pitchAndFilterEnvShared || v.filterEnvActive {
		if doneCB := v.filterEnv.SetEnvelopePosition(pos); doneCB != nil {
			doneCB(v)
		}
	}
}

func (v itVoice[TPeriod]) GetFilterEnvelopePosition() int {
	return v.filterEnv.GetEnvelopePosition()
}
