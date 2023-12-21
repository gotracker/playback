package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// PitchEnvelopeOff defines a panning envelope: off effect
type PitchEnvelopeOff[TPeriod period.Period] DataEffect // 'S7B'

// Start triggers on the first tick, but before the Tick() function is called
func (e PitchEnvelopeOff[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	cs.SetPitchEnvelopeEnable(false)
	return nil
}

func (e PitchEnvelopeOff[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
