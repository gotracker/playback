package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// PortaDown defines a portamento down effect
type PortaDown[TPeriod period.Period] DataEffect // '2'

// Start triggers on the first tick, but before the Tick() function is called
func (e PortaDown[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e PortaDown[TPeriod]) Tick(cs playback.Channel[TPeriod, Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	xx := mem.PortaDown(DataEffect(e))

	if currentTick == 0 {
		return nil
	}

	return doPortaDown(cs, float32(xx), 4)
}

func (e PortaDown[TPeriod]) String() string {
	return fmt.Sprintf("2%0.2x", DataEffect(e))
}
