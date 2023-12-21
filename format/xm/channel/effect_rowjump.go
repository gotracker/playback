package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
)

// RowJump defines a row jump effect
type RowJump[TPeriod period.Period] DataEffect // 'D'

// Start triggers on the first tick, but before the Tick() function is called
func (e RowJump[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Stop is called on the last tick of the row, but after the Tick() function is called
func (e RowJump[TPeriod]) Stop(cs playback.Channel[TPeriod, Memory], p playback.Playback, lastTick int) error {
	xy := DataEffect(e)
	x := xy >> 4
	y := xy & 0x0f
	row := index.Row(x*10 + y)
	if err := p.BreakOrder(); err != nil {
		return err
	}
	return p.SetNextRow(row)
}

func (e RowJump[TPeriod]) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}
