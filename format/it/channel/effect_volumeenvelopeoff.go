package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// VolumeEnvelopeOff defines a volume envelope: off effect
type VolumeEnvelopeOff[TPeriod period.Period] DataEffect // 'S77'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeEnvelopeOff[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	cs.SetVolumeEnvelopeEnable(false)
	return nil
}

func (e VolumeEnvelopeOff[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
