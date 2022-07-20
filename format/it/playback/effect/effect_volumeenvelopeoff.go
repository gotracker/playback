package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/layout/channel"
)

// VolumeEnvelopeOff defines a volume envelope: off effect
type VolumeEnvelopeOff channel.DataEffect // 'S77'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeEnvelopeOff) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	cs.SetVolumeEnvelopeEnable(false)
	return nil
}

func (e VolumeEnvelopeOff) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
