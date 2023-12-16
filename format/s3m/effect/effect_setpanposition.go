package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
)

// SetPanPosition defines a set pan position effect
type SetPanPosition ChannelCommand // 'S8x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetPanPosition) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := uint8(e) & 0xf

	cs.SetPan(s3mPanning.PanningFromS3M(x))
	return nil
}

func (e SetPanPosition) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
