package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
)

// VolumeEnvelopeOn defines a volume envelope: on effect
type VolumeEnvelopeOn channel.DataEffect // 'S78'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeEnvelopeOn) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	cs.SetVolumeEnvelopeEnable(true)
	return nil
}

func (e VolumeEnvelopeOn) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
