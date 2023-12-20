package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// FineVibrato defines an fine vibrato effect
type FineVibrato[TPeriod period.Period] channel.DataEffect // 'U'

// Start triggers on the first tick, but before the Tick() function is called
func (e FineVibrato[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e FineVibrato[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.Vibrato(channel.DataEffect(e))
	if currentTick != 0 {
		return doVibrato(cs, currentTick, x, y, 1)
	}
	return nil
}

func (e FineVibrato[TPeriod]) String() string {
	return fmt.Sprintf("U%0.2x", channel.DataEffect(e))
}
