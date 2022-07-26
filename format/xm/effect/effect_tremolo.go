package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// Tremolo defines a tremolo effect
type Tremolo channel.DataEffect // '7'

// Start triggers on the first tick, but before the Tick() function is called
func (e Tremolo) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e Tremolo) Tick(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.Tremolo(channel.DataEffect(e))
	// NOTE: JBC - XM updates on tick 0, but MOD does not.
	// Just have to eat this incompatibility, I guess...
	return doTremolo(cs, currentTick, x, y, 4)
}

func (e Tremolo) String() string {
	return fmt.Sprintf("7%0.2x", channel.DataEffect(e))
}
