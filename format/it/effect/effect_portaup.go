package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// PortaUp defines a portamento up effect
type PortaUp[TPeriod period.Period] channel.DataEffect // 'F'

// Start triggers on the first tick, but before the Tick() function is called
func (e PortaUp[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e PortaUp[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	xx := mem.PortaUp(channel.DataEffect(e))

	return doPortaUp(cs, float32(xx), 4)
}

func (e PortaUp[TPeriod]) String() string {
	return fmt.Sprintf("F%0.2x", channel.DataEffect(e))
}
