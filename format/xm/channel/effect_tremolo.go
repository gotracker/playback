package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// Tremolo defines a tremolo effect
type Tremolo[TPeriod period.Period] DataEffect // '7'

// Start triggers on the first tick, but before the Tick() function is called
func (e Tremolo[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e Tremolo[TPeriod]) Tick(cs playback.Channel[TPeriod, Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.Tremolo(DataEffect(e))
	// NOTE: JBC - XM updates on tick 0, but MOD does not.
	// Just have to eat this incompatibility, I guess...
	return doTremolo(cs, currentTick, x, y, 4)
}

func (e Tremolo[TPeriod]) String() string {
	return fmt.Sprintf("7%0.2x", DataEffect(e))
}
