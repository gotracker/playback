package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/layout/channel"
)

// ExtraFinePortaDown defines an extra-fine portamento down effect
type ExtraFinePortaDown channel.DataEffect // 'EEx'

// Start triggers on the first tick, but before the Tick() function is called
func (e ExtraFinePortaDown) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	mem := cs.GetMemory()
	y := mem.PortaDown(channel.DataEffect(e)) & 0x0F

	return doPortaDown(cs, float32(y), 1, mem.Shared.LinearFreqSlides)
}

func (e ExtraFinePortaDown) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
