package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// Tremolo defines a tremolo effect
type Tremolo[TPeriod period.Period] channel.DataEffect // 'R'

// Start triggers on the first tick, but before the Tick() function is called
func (e Tremolo[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e Tremolo[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.Tremolo(channel.DataEffect(e))
	// NOTE: JBC - IT dos not update on tick 0, but MOD does.
	// Maybe need to add a flag for converted MOD backward compatibility?
	if currentTick != 0 {
		return doTremolo(cs, currentTick, x, y, 4)
	}
	return nil
}

func (e Tremolo[TPeriod]) String() string {
	return fmt.Sprintf("R%0.2x", channel.DataEffect(e))
}
