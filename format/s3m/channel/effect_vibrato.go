package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// Vibrato defines a vibrato effect
type Vibrato ChannelCommand // 'H'

// Start triggers on the first tick, but before the Tick() function is called
func (e Vibrato) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e Vibrato) Tick(cs S3MChannel, p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.Vibrato(DataEffect(e))
	// NOTE: JBC - S3M dos not update on tick 0, but MOD does.
	if currentTick != 0 || mem.Shared.ModCompatibility {
		return doVibrato(cs, currentTick, x, y, 4)
	}
	return nil
}

func (e Vibrato) String() string {
	return fmt.Sprintf("H%0.2x", DataEffect(e))
}