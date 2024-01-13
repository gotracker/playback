package voice

import (
	"github.com/gotracker/playback/period"
)

// == PitchEnveloper ==

func (v *xmVoice[TPeriod]) EnablePitchEnvelope(enabled bool) {
}

func (v xmVoice[TPeriod]) IsPitchEnvelopeEnabled() bool {
	return false
}

func (v xmVoice[TPeriod]) GetCurrentPitchEnvelope() period.Delta {
	return 0
}

func (v *xmVoice[TPeriod]) SetPitchEnvelopePosition(pos int) {
}

func (v xmVoice[TPeriod]) GetPitchEnvelopePosition() int {
	return 0
}
