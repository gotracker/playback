package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
)

// ExtraFinePortaUp defines an extra-fine portamento up effect
type ExtraFinePortaUp ChannelCommand // 'FEx'

// Start triggers on the first tick, but before the Tick() function is called
func (e ExtraFinePortaUp) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	y := channel.DataEffect(e) & 0x0F

	return doPortaUp(cs, float32(y), 1)
}

func (e ExtraFinePortaUp) String() string {
	return fmt.Sprintf("F%0.2x", channel.DataEffect(e))
}
