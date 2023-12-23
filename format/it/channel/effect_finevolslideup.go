package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// FineVolumeSlideUp defines a fine volume slide up effect
type FineVolumeSlideUp[TPeriod period.Period] DataEffect // 'D'

// Start triggers on the first tick, but before the Tick() function is called
func (e FineVolumeSlideUp[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e FineVolumeSlideUp[TPeriod]) Tick(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, _ := mem.VolumeSlide(DataEffect(e))

	if x != 0x0F && currentTick == 0 {
		return doVolSlide(cs, float32(x), 1.0)
	}
	return nil
}

func (e FineVolumeSlideUp[TPeriod]) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}

//====================================================

// VolChanFineVolumeSlideUp defines a fine volume slide up effect (from the volume channel)
type VolChanFineVolumeSlideUp[TPeriod period.Period] DataEffect // 'd'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolChanFineVolumeSlideUp[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	mem := cs.GetMemory()
	x := mem.VolChanVolumeSlide(DataEffect(e))

	return doVolSlide(cs, float32(x), 1.0)
}

func (e VolChanFineVolumeSlideUp[TPeriod]) String() string {
	return fmt.Sprintf("d%xF", DataEffect(e))
}
