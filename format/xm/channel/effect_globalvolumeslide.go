package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// GlobalVolumeSlide defines a global volume slide effect
type GlobalVolumeSlide[TPeriod period.Period] DataEffect // 'H'

// Start triggers on the first tick, but before the Tick() function is called
func (e GlobalVolumeSlide[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e GlobalVolumeSlide[TPeriod]) Tick(cs playback.Channel[TPeriod, Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.GlobalVolumeSlide(DataEffect(e))

	if currentTick == 0 {
		return nil
	}

	m := p.(XM)

	if x == 0 {
		// global vol slide down
		return doGlobalVolSlide(m, -float32(y), 1.0)
	} else if y == 0 {
		// global vol slide up
		return doGlobalVolSlide(m, float32(y), 1.0)
	}
	return nil
}

func (e GlobalVolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("H%0.2x", DataEffect(e))
}
