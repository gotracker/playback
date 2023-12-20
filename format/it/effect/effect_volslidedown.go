package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// VolumeSlideDown defines a volume slide down effect
type VolumeSlideDown[TPeriod period.Period] channel.DataEffect // 'D'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeSlideDown[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e VolumeSlideDown[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	_, y := mem.VolumeSlide(channel.DataEffect(e))

	return doVolSlide(cs, -float32(y), 1.0)
}

func (e VolumeSlideDown[TPeriod]) String() string {
	return fmt.Sprintf("D%0.2x", channel.DataEffect(e))
}

//====================================================

// VolChanVolumeSlideDown defines a volume slide down effect (from the volume channel)
type VolChanVolumeSlideDown[TPeriod period.Period] channel.DataEffect // 'd'

// Tick is called on every tick
func (e VolChanVolumeSlideDown[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	y := mem.VolChanVolumeSlide(channel.DataEffect(e))

	return doVolSlide(cs, -float32(y), 1.0)
}

func (e VolChanVolumeSlideDown[TPeriod]) String() string {
	return fmt.Sprintf("d0%x", channel.DataEffect(e))
}
