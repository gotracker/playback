package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
)

// ExtraFinePortaDown defines an extra-fine portamento down effect
type ExtraFinePortaDown ChannelCommand // 'EEx'

// Start triggers on the first tick, but before the Tick() function is called
func (e ExtraFinePortaDown) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	y := channel.DataEffect(e) & 0x0F

	return doPortaDown(cs, float32(y), 1)
}

func (e ExtraFinePortaDown) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
