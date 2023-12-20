package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// ExtraFinePortaUp defines an extra-fine portamento up effect
type ExtraFinePortaUp[TPeriod period.Period] channel.DataEffect // 'FEx'

// Start triggers on the first tick, but before the Tick() function is called
func (e ExtraFinePortaUp[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	mem := cs.GetMemory()
	y := mem.PortaUp(channel.DataEffect(e)) & 0x0F

	return doPortaUp(cs, float32(y), 1)
}

func (e ExtraFinePortaUp[TPeriod]) String() string {
	return fmt.Sprintf("F%0.2x", channel.DataEffect(e))
}
