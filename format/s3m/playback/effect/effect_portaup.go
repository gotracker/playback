package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/layout/channel"
)

// PortaUp defines a portamento up effect
type PortaUp ChannelCommand // 'F'

// Start triggers on the first tick, but before the Tick() function is called
func (e PortaUp) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e PortaUp) Tick(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	xx := mem.LastNonZero(channel.DataEffect(e))

	if currentTick != 0 {
		return doPortaUp(cs, float32(xx), 4)
	}
	return nil
}

func (e PortaUp) String() string {
	return fmt.Sprintf("F%0.2x", channel.DataEffect(e))
}
