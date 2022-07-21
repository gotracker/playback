package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
)

// PortaDown defines a portamento down effect
type PortaDown channel.DataEffect // 'E'

// Start triggers on the first tick, but before the Tick() function is called
func (e PortaDown) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e PortaDown) Tick(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	xx := mem.PortaDown(channel.DataEffect(e))

	return doPortaDown(cs, float32(xx), 4, mem.Shared.LinearFreqSlides)
}

func (e PortaDown) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
