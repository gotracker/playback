package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// VolumeEnvelopeOn defines a volume envelope: on effect
type VolumeEnvelopeOn[TPeriod period.Period] channel.DataEffect // 'S78'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeEnvelopeOn[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	cs.SetVolumeEnvelopeEnable(true)
	return nil
}

func (e VolumeEnvelopeOn[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
