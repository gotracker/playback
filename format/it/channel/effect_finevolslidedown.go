package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// FineVolumeSlideDown defines a fine volume slide down effect
type FineVolumeSlideDown[TPeriod period.Period] DataEffect // 'D'

// Start triggers on the first tick, but before the Tick() function is called
func (e FineVolumeSlideDown[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e FineVolumeSlideDown[TPeriod]) Tick(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	_, y := mem.VolumeSlide(DataEffect(e))

	if y != 0x0F && currentTick == 0 {
		return doVolSlide(cs, -float32(y), 1.0)
	}
	return nil
}

func (e FineVolumeSlideDown[TPeriod]) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}

//====================================================

// VolChanFineVolumeSlideDown defines a fine volume slide down effect (from the volume channel)
type VolChanFineVolumeSlideDown[TPeriod period.Period] DataEffect // 'd'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolChanFineVolumeSlideDown[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	mem := cs.GetMemory()
	y := mem.VolChanVolumeSlide(DataEffect(e))

	return doVolSlide(cs, -float32(y), 1.0)
}

func (e VolChanFineVolumeSlideDown[TPeriod]) String() string {
	return fmt.Sprintf("dF%x", DataEffect(e))
}
