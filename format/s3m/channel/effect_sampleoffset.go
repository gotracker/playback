package channel

import (
	"fmt"

	"github.com/gotracker/gomixing/sampling"

	"github.com/gotracker/playback"
)

// SampleOffset defines a sample offset effect
type SampleOffset ChannelCommand // 'O'

// Start triggers on the first tick, but before the Tick() function is called
func (e SampleOffset) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	mem := cs.GetMemory()
	xx := mem.SampleOffset(DataEffect(e))
	cs.SetTargetPos(sampling.Pos{Pos: int(xx) * 0x100})
	return nil
}

func (e SampleOffset) String() string {
	return fmt.Sprintf("O%0.2x", DataEffect(e))
}
