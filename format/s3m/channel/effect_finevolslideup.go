package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// FineVolumeSlideUp defines a fine volume slide up effect
type FineVolumeSlideUp ChannelCommand // 'DxF'

// Start triggers on the first tick, but before the Tick() function is called
func (e FineVolumeSlideUp) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e FineVolumeSlideUp) Tick(cs S3MChannel, p playback.Playback, currentTick int) error {
	x := DataEffect(e) >> 4

	if x != 0x0F && currentTick == 0 {
		return doVolSlide(cs, float32(x), 1.0)
	}
	return nil
}

func (e FineVolumeSlideUp) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}
