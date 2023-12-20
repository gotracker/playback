package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// VolumeEnvelopeOff defines a volume envelope: off effect
type VolumeEnvelopeOff[TPeriod period.Period] channel.DataEffect // 'S77'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeEnvelopeOff[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	cs.SetVolumeEnvelopeEnable(false)
	return nil
}

func (e VolumeEnvelopeOff[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
