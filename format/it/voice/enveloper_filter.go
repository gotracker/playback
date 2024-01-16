package voice

// == FilterEnveloper ==

func (v *itVoice[TPeriod]) EnableFilterEnvelope(enabled bool) error {
	if !v.pitchAndFilterEnvShared {
		return v.filterEnv.SetEnabled(enabled)
	}

	// shared filter/pitch envelope
	if !v.filterEnvActive {
		return nil
	}

	return v.filterEnv.SetEnabled(enabled)
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

func (v *itVoice[TPeriod]) SetFilterEnvelopePosition(pos int) error {
	if !v.pitchAndFilterEnvShared || v.filterEnvActive {
		doneCB, err := v.filterEnv.SetEnvelopePosition(pos)
		if err != nil {
			return err
		}
		if doneCB != nil {
			doneCB(v)
		}
	}
	return nil
}

func (v itVoice[TPeriod]) GetFilterEnvelopePosition() int {
	return v.filterEnv.GetEnvelopePosition()
}
