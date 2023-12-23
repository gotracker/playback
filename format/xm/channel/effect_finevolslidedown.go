package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// FineVolumeSlideDown defines a volume slide effect
type FineVolumeSlideDown[TPeriod period.Period] DataEffect // 'EAx'

// Start triggers on the first tick, but before the Tick() function is called
func (e FineVolumeSlideDown[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	mem := cs.GetMemory()
	xy := mem.FineVolumeSlideDown(DataEffect(e))
	y := DataEffect(xy & 0x0F)

	return doVolSlide(cs, -float32(y), 1.0)
}

func (e FineVolumeSlideDown[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}
