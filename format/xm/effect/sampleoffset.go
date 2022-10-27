package effect

import (
	"fmt"

	"github.com/gotracker/gomixing/sampling"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// SampleOffset defines a sample offset effect
type SampleOffset channel.DataEffect // '9'

// Start triggers on the first tick, but before the Tick() function is called
func (e SampleOffset) Start(cs *channel.State, p playback.Playback) error {
	cs.ResetRetriggerCount()
	mem := cs.GetMemory()
	xx := mem.SampleOffset(channel.DataEffect(e))
	cs.SetTargetPos(sampling.Pos{Pos: int(xx) * 0x100})
	return nil
}

func (e SampleOffset) String() string {
	return fmt.Sprintf("9%0.2x", channel.DataEffect(e))
}
