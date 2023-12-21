package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// Tremor defines a tremor effect
type Tremor ChannelCommand // 'I'

// Start triggers on the first tick, but before the Tick() function is called
func (e Tremor) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e Tremor) Tick(cs S3MChannel, p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.LastNonZeroXY(DataEffect(e))
	return doTremor(cs, currentTick, int(x)+1, int(y)+1)
}

func (e Tremor) String() string {
	return fmt.Sprintf("I%0.2x", DataEffect(e))
}
