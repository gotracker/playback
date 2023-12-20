package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/period"
)

// FineVolumeSlideDown defines a volume slide effect
type FineVolumeSlideDown[TPeriod period.Period] channel.DataEffect // 'EAx'

// Start triggers on the first tick, but before the Tick() function is called
func (e FineVolumeSlideDown[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	mem := cs.GetMemory()
	xy := mem.FineVolumeSlideDown(channel.DataEffect(e))
	y := channel.DataEffect(xy & 0x0F)

	return doVolSlide(cs, -float32(y), 1.0)
}

func (e FineVolumeSlideDown[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
