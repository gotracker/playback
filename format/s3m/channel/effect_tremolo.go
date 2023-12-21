package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// Tremolo defines a tremolo effect
type Tremolo ChannelCommand // 'R'

// Start triggers on the first tick, but before the Tick() function is called
func (e Tremolo) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e Tremolo) Tick(cs S3MChannel, p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.Tremolo(DataEffect(e))
	// NOTE: JBC - S3M does not update on tick 0, but MOD does.
	if currentTick != 0 || mem.Shared.ModCompatibility {
		return doTremolo(cs, currentTick, DataEffect(x), DataEffect(y), 4)
	}
	return nil
}

func (e Tremolo) String() string {
	return fmt.Sprintf("R%0.2x", DataEffect(e))
}
