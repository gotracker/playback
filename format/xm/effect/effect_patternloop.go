package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/period"
)

// PatternLoop defines a pattern loop effect
type PatternLoop[TPeriod period.Period] channel.DataEffect // 'E6x'

// Start triggers on the first tick, but before the Tick() function is called
func (e PatternLoop[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := channel.DataEffect(e) & 0xF

	mem := cs.GetMemory()
	pl := mem.GetPatternLoop()
	if x == 0 {
		// set loop
		pl.Start = p.GetCurrentRow()
	} else {
		if !pl.Enabled {
			pl.Enabled = true
			pl.Total = uint8(x)
			pl.End = p.GetCurrentRow()
			pl.Count = 0
		}
		if row, ok := pl.ContinueLoop(p.GetCurrentRow()); ok {
			return p.SetNextRowWithBacktrack(row, true)
		}
	}
	return nil
}

func (e PatternLoop[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
