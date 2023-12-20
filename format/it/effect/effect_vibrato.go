package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// Vibrato defines a vibrato effect
type Vibrato[TPeriod period.Period] channel.DataEffect // 'H'

// Start triggers on the first tick, but before the Tick() function is called
func (e Vibrato[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e Vibrato[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.Vibrato(channel.DataEffect(e))
	if mem.Shared.OldEffectMode {
		if currentTick != 0 {
			return doVibrato(cs, currentTick, x, y, 8)
		}
	} else {
		return doVibrato(cs, currentTick, x, y, 4)
	}
	return nil
}

func (e Vibrato[TPeriod]) String() string {
	return fmt.Sprintf("H%0.2x", channel.DataEffect(e))
}
