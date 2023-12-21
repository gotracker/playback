package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
)

// SetPanPosition defines a set pan position effect
type SetPanPosition ChannelCommand // 'S8x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetPanPosition) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := uint8(e) & 0xf

	cs.SetPan(s3mPanning.PanningFromS3M(x))
	return nil
}

func (e SetPanPosition) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
