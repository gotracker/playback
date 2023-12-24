package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
)

// StereoControl defines a set stereo control effect
type StereoControl ChannelCommand // 'SAx'

// Start triggers on the first tick, but before the Tick() function is called
func (e StereoControl) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := uint8(e) & 0xf

	active := cs.GetActiveState()
	if x > 7 {
		active.Pan = s3mPanning.PanningFromS3M(x - 8)
	} else {
		active.Pan = s3mPanning.PanningFromS3M(x + 8)
	}
	return nil
}

func (e StereoControl) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
