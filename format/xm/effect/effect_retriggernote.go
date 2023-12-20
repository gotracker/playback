package effect

import (
	"fmt"

	"github.com/gotracker/gomixing/sampling"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/period"
)

// RetriggerNote defines a retriggering effect
type RetriggerNote[TPeriod period.Period] channel.DataEffect // 'E9x'

// Start triggers on the first tick, but before the Tick() function is called
func (e RetriggerNote[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e RetriggerNote[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
	y := channel.DataEffect(e) & 0x0F
	if y == 0 {
		return nil
	}

	rt := cs.GetRetriggerCount() + 1
	cs.SetRetriggerCount(rt)
	if channel.DataEffect(rt) >= y {
		cs.SetPos(sampling.Pos{})
		cs.ResetRetriggerCount()
	}
	return nil
}

func (e RetriggerNote[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
