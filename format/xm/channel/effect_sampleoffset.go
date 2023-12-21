package channel

import (
	"fmt"

	"github.com/gotracker/gomixing/sampling"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// SampleOffset defines a sample offset effect
type SampleOffset[TPeriod period.Period] DataEffect // '9'

// Start triggers on the first tick, but before the Tick() function is called
func (e SampleOffset[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	mem := cs.GetMemory()
	xx := mem.SampleOffset(DataEffect(e))
	cs.SetTargetPos(sampling.Pos{Pos: int(xx) * 0x100})
	return nil
}

func (e SampleOffset[TPeriod]) String() string {
	return fmt.Sprintf("9%0.2x", DataEffect(e))
}
