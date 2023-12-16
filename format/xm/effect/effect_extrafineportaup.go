package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// ExtraFinePortaUp defines an extra-fine portamento up effect
type ExtraFinePortaUp channel.DataEffect // 'X1x'

// Start triggers on the first tick, but before the Tick() function is called
func (e ExtraFinePortaUp) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	mem := cs.GetMemory()
	xx := mem.ExtraFinePortaUp(channel.DataEffect(e))
	y := xx & 0x0F

	return doPortaUp(cs, float32(y), 1, mem.Shared.LinearFreqSlides)
}

func (e ExtraFinePortaUp) String() string {
	return fmt.Sprintf("F%0.2x", channel.DataEffect(e))
}
