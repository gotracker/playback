package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// VolumeSlideUp defines a volume slide up effect
type VolumeSlideUp[TPeriod period.Period] channel.DataEffect // 'D'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolumeSlideUp[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e VolumeSlideUp[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, _ := mem.VolumeSlide(channel.DataEffect(e))

	return doVolSlide(cs, float32(x), 1.0)
}

func (e VolumeSlideUp[TPeriod]) String() string {
	return fmt.Sprintf("D%0.2x", channel.DataEffect(e))
}

//====================================================

// VolChanVolumeSlideUp defines a volume slide up effect (from the volume channel)
type VolChanVolumeSlideUp[TPeriod period.Period] channel.DataEffect // 'd'

// Tick is called on every tick
func (e VolChanVolumeSlideUp[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x := mem.VolChanVolumeSlide(channel.DataEffect(e))

	return doVolSlide(cs, float32(x), 1.0)
}

func (e VolChanVolumeSlideUp[TPeriod]) String() string {
	return fmt.Sprintf("d%x0", channel.DataEffect(e))
}
