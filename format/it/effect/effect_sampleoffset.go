package effect

import (
	"fmt"

	"github.com/gotracker/gomixing/sampling"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
)

// SampleOffset defines a sample offset effect
type SampleOffset channel.DataEffect // 'O'

// Start triggers on the first tick, but before the Tick() function is called
func (e SampleOffset) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	mem := cs.GetMemory()
	xx := mem.SampleOffset(channel.DataEffect(e))

	pos := sampling.Pos{Pos: mem.HighOffset + int(xx)*0x100}
	if mem.Shared.OldEffectMode {
		if inst := cs.GetInstrument(); inst != nil && inst.GetLength().Pos < pos.Pos {
			cs.SetTargetPos(pos)
		}
	} else {
		cs.SetTargetPos(pos)
	}
	return nil
}

func (e SampleOffset) String() string {
	return fmt.Sprintf("O%0.2x", channel.DataEffect(e))
}
