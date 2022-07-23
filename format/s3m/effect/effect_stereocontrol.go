package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
)

// StereoControl defines a set stereo control effect
type StereoControl ChannelCommand // 'SAx'

// Start triggers on the first tick, but before the Tick() function is called
func (e StereoControl) Start(cs *channel.State, p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := uint8(e) & 0xf

	if x > 7 {
		cs.SetPan(s3mPanning.PanningFromS3M(x - 8))
	} else {
		cs.SetPan(s3mPanning.PanningFromS3M(x + 8))
	}
	return nil
}

func (e StereoControl) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
