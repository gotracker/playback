package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// PatternLoop defines a pattern loop effect
type PatternLoop ChannelCommand // 'SBx'

// Start triggers on the first tick, but before the Tick() function is called
func (e PatternLoop) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Stop is called on the last tick of the row, but after the Tick() function is called
func (e PatternLoop) Stop(cs S3MChannel, p playback.Playback, lastTick int) error {
	x := DataEffect(e) & 0xF

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

func (e PatternLoop) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
