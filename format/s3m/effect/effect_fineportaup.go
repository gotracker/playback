package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
)

// FinePortaUp defines an fine portamento up effect
type FinePortaUp ChannelCommand // 'FFx'

// Start triggers on the first tick, but before the Tick() function is called
func (e FinePortaUp) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	y := channel.DataEffect(e) & 0x0F

	return doPortaUp(cs, float32(y), 4)
}

func (e FinePortaUp) String() string {
	return fmt.Sprintf("F%0.2x", channel.DataEffect(e))
}
