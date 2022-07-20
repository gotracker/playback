package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
)

// PitchEnvelopeOn defines a panning envelope: on effect
type PitchEnvelopeOn channel.DataEffect // 'S7C'

// Start triggers on the first tick, but before the Tick() function is called
func (e PitchEnvelopeOn) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	cs.SetPitchEnvelopeEnable(true)
	return nil
}

func (e PitchEnvelopeOn) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
