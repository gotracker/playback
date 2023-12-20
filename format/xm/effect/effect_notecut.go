package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/period"
)

// NoteCut defines a note cut effect
type NoteCut[TPeriod period.Period] channel.DataEffect // 'ECx'

// Start triggers on the first tick, but before the Tick() function is called
func (e NoteCut[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e NoteCut[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
	x := channel.DataEffect(e) & 0xf

	if x != 0 && currentTick == int(x) {
		cs.NoteCut()
	}
	return nil
}

func (e NoteCut[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
