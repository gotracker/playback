package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// PortaDown defines a portamento down effect
type PortaDown ChannelCommand // 'E'

// Start triggers on the first tick, but before the Tick() function is called
func (e PortaDown) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e PortaDown) Tick(cs S3MChannel, p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	xx := mem.LastNonZero(DataEffect(e))

	if currentTick != 0 {
		return doPortaDown(cs, float32(xx), 4)
	}
	return nil
}

func (e PortaDown) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}