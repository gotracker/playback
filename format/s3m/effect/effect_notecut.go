package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
)

// NoteCut defines a note cut effect
type NoteCut ChannelCommand // 'SCx'

// Start triggers on the first tick, but before the Tick() function is called
func (e NoteCut) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e NoteCut) Tick(cs playback.Channel[channel.Memory], p playback.Playback, currentTick int) error {
	x := channel.DataEffect(e) & 0xf

	if x != 0 && currentTick == int(x) {
		cs.FreezePlayback()
	}
	return nil
}

func (e NoteCut) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
