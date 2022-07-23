package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
)

// Vibrato defines a vibrato effect
type Vibrato ChannelCommand // 'H'

// Start triggers on the first tick, but before the Tick() function is called
func (e Vibrato) Start(cs *channel.State, p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e Vibrato) Tick(cs *channel.State, p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.Vibrato(channel.DataEffect(e))
	// NOTE: JBC - S3M dos not update on tick 0, but MOD does.
	if currentTick != 0 || mem.Shared.ModCompatibility {
		return doVibrato(cs, currentTick, x, y, 4)
	}
	return nil
}

func (e Vibrato) String() string {
	return fmt.Sprintf("H%0.2x", channel.DataEffect(e))
}
