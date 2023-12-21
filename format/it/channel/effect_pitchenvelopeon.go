package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// PitchEnvelopeOn defines a panning envelope: on effect
type PitchEnvelopeOn[TPeriod period.Period] DataEffect // 'S7C'

// Start triggers on the first tick, but before the Tick() function is called
func (e PitchEnvelopeOn[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	cs.SetPitchEnvelopeEnable(true)
	return nil
}

func (e PitchEnvelopeOn[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
