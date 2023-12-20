package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/period"
)

// VolumeSlide defines a volume slide effect
type VolumeSlide[TPeriod period.Period] channel.DataEffect // 'A'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeSlide[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e VolumeSlide[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
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

func (e VolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("A%0.2x", channel.DataEffect(e))
}
