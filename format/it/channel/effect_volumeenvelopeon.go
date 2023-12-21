package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// VolumeEnvelopeOn defines a volume envelope: on effect
type VolumeEnvelopeOn[TPeriod period.Period] DataEffect // 'S78'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeEnvelopeOn[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	cs.SetVolumeEnvelopeEnable(true)
	return nil
}

func (e VolumeEnvelopeOn[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
