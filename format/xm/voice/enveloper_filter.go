package voice

// == FilterEnveloper ==

func (v *xmVoice[TPeriod]) EnableFilterEnvelope(enabled bool) {
}

func (v xmVoice[TPeriod]) IsFilterEnvelopeEnabled() bool {
	return false
}

func (v xmVoice[TPeriod]) GetCurrentFilterEnvelope() uint8 {
	return 0
}

func (v *xmVoice[TPeriod]) SetFilterEnvelopePosition(pos int) {
}

func (v xmVoice[TPeriod]) GetFilterEnvelopePosition() int {
	return 0
}
