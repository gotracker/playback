package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// Vibrato defines a vibrato effect
type Vibrato[TPeriod period.Period] DataEffect // '4'

// Start triggers on the first tick, but before the Tick() function is called
func (e Vibrato[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e Vibrato[TPeriod]) Tick(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.Vibrato(DataEffect(e))
	// NOTE: JBC - XM updates on tick 0, but MOD does not.
	// Just have to eat this incompatibility, I guess...
	return doVibrato(cs, currentTick, x, y, 4)
}

func (e Vibrato[TPeriod]) String() string {
	return fmt.Sprintf("4%0.2x", DataEffect(e))
}
