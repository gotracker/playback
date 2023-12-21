package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// NoteCut defines a note cut effect
type NoteCut ChannelCommand // 'SCx'

// Start triggers on the first tick, but before the Tick() function is called
func (e NoteCut) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e NoteCut) Tick(cs S3MChannel, p playback.Playback, currentTick int) error {
	x := DataEffect(e) & 0xf

	if x != 0 && currentTick == int(x) {
		cs.FreezePlayback()
	}
	return nil
}

func (e NoteCut) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}