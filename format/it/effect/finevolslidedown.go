package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
)

// FineVolumeSlideDown defines a fine volume slide down effect
type FineVolumeSlideDown channel.DataEffect // 'D'

// Start triggers on the first tick, but before the Tick() function is called
func (e FineVolumeSlideDown) Start(cs *channel.State, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e FineVolumeSlideDown) Tick(cs *channel.State, p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	_, y := mem.VolumeSlide(channel.DataEffect(e))

	if y != 0x0F && currentTick == 0 {
		return doVolSlide(cs, -float32(y), 1.0)
	}
	return nil
}

func (e FineVolumeSlideDown) String() string {
	return fmt.Sprintf("D%0.2x", channel.DataEffect(e))
}

//====================================================

// VolChanFineVolumeSlideDown defines a fine volume slide down effect (from the volume channel)
type VolChanFineVolumeSlideDown channel.DataEffect // 'd'

// Start triggers on the first tick, but before the Tick() function is called
func (e VolChanFineVolumeSlideDown) Start(cs *channel.State, p playback.Playback) error {
	mem := cs.GetMemory()
	y := mem.VolChanVolumeSlide(channel.DataEffect(e))

	return doVolSlide(cs, -float32(y), 1.0)
}

func (e VolChanFineVolumeSlideDown) String() string {
	return fmt.Sprintf("dF%x", channel.DataEffect(e))
}
