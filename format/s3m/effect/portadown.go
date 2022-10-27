package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
)

// PortaDown defines a portamento down effect
type PortaDown ChannelCommand // 'E'

// Start triggers on the first tick, but before the Tick() function is called
func (e PortaDown) Start(cs *channel.State, p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e PortaDown) Tick(cs *channel.State, p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	xx := mem.LastNonZero(channel.DataEffect(e))

	if currentTick != 0 {
		return doPortaDown(cs, float32(xx), 4)
	}
	return nil
}

func (e PortaDown) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
