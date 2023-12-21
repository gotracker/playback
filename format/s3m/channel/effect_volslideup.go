package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// VolumeSlideUp defines a volume slide up effect
type VolumeSlideUp ChannelCommand // 'Dx0'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeSlideUp) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e VolumeSlideUp) Tick(cs S3MChannel, p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x := DataEffect(e) >> 4

	if mem.Shared.VolSlideEveryFrame || currentTick != 0 {
		return doVolSlide(cs, float32(x), 1.0)
	}
	return nil
}

func (e VolumeSlideUp) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}
