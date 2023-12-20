package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
)

// RowJump defines a row jump effect
type RowJump[TPeriod period.Period] channel.DataEffect // 'C'

// Start triggers on the first tick, but before the Tick() function is called
func (e RowJump[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Stop is called on the last tick of the row, but after the Tick() function is called
func (e RowJump[TPeriod]) Stop(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, lastTick int) error {
	r := channel.DataEffect(e)
	rowIdx := index.Row(r)
	return p.SetNextRow(rowIdx)
}

func (e RowJump[TPeriod]) String() string {
	return fmt.Sprintf("C%0.2x", channel.DataEffect(e))
}
