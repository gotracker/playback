package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
)

// SurroundOn defines a set surround on effect
type SurroundOn ChannelCommand // 'S91'

// Start triggers on the first tick, but before the Tick() function is called
func (e SurroundOn) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()

	// TODO: support for center rear panning
	cs.GetActiveState().Pan = s3mPanning.DefaultPanning
	return nil
}

func (e SurroundOn) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
