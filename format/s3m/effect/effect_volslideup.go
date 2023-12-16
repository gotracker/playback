package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
)

// VolumeSlideUp defines a volume slide up effect
type VolumeSlideUp ChannelCommand // 'Dx0'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeSlideUp) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e VolumeSlideUp) Tick(cs playback.Channel[channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x := channel.DataEffect(e) >> 4

	if mem.Shared.VolSlideEveryFrame || currentTick != 0 {
		return doVolSlide(cs, float32(x), 1.0)
	}
	return nil
}

func (e VolumeSlideUp) String() string {
	return fmt.Sprintf("D%0.2x", channel.DataEffect(e))
}
