package effect

import (
	"fmt"

	"github.com/gotracker/gomixing/sampling"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// RetriggerNote defines a retriggering effect
type RetriggerNote channel.DataEffect // 'E9x'

// Start triggers on the first tick, but before the Tick() function is called
func (e RetriggerNote) Start(cs *channel.State, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e RetriggerNote) Tick(cs *channel.State, p playback.Playback, currentTick int) error {
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

func (e RetriggerNote) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
