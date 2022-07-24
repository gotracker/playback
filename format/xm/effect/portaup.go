package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// PortaUp defines a portamento up effect
type PortaUp channel.DataEffect // '1'

// Start triggers on the first tick, but before the Tick() function is called
func (e PortaUp) Start(cs *channel.State, p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e PortaUp) Tick(cs *channel.State, p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	xx := mem.PortaUp(channel.DataEffect(e))

	if currentTick == 0 {
		return nil
	}

	return doPortaUp(cs, float32(xx), 4)
}

func (e PortaUp) String() string {
	return fmt.Sprintf("1%0.2x", channel.DataEffect(e))
}
