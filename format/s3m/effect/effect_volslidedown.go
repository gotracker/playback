package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
)

// VolumeSlideDown defines a volume slide down effect
type VolumeSlideDown ChannelCommand // 'D0y'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeSlideDown) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e VolumeSlideDown) Tick(cs playback.Channel[channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	y := channel.DataEffect(e) & 0x0F

	if mem.Shared.VolSlideEveryFrame || currentTick != 0 {
		return doVolSlide(cs, -float32(y), 1.0)
	}
	return nil
}

func (e VolumeSlideDown) String() string {
	return fmt.Sprintf("D%0.2x", channel.DataEffect(e))
}
