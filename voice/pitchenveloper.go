package voice

import (
	"github.com/gotracker/playback/period"
)

// PitchEnveloper is a pitch envelope interface
type PitchEnveloper[TPeriod period.Period] interface {
	EnablePitchEnvelope(enabled bool)
	IsPitchEnvelopeEnabled() bool
	GetCurrentPitchEnvelope() period.Delta
	SetPitchEnvelopePosition(pos int)
}
