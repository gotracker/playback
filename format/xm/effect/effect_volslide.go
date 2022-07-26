package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// VolumeSlide defines a volume slide effect
type VolumeSlide channel.DataEffect // 'A'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeSlide) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e VolumeSlide) Tick(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.VolumeSlide(channel.DataEffect(e))

	if currentTick == 0 {
		return nil
	}

	if x == 0 {
		// vol slide down
		return doVolSlide(cs, -float32(y), 1.0)
	} else if y == 0 {
		// vol slide up
		return doVolSlide(cs, float32(y), 1.0)
	}
	return nil
}

func (e VolumeSlide) String() string {
	return fmt.Sprintf("A%0.2x", channel.DataEffect(e))
}
