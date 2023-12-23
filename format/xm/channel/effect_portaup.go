package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// PortaUp defines a portamento up effect
type PortaUp[TPeriod period.Period] DataEffect // '1'

// Start triggers on the first tick, but before the Tick() function is called
func (e PortaUp[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

// Tick is called on every tick
func (e PortaUp[TPeriod]) Tick(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	xx := mem.PortaUp(DataEffect(e))

	if currentTick == 0 {
		return nil
	}

	return doPortaUp(cs, float32(xx), 4)
}

func (e PortaUp[TPeriod]) String() string {
	return fmt.Sprintf("1%0.2x", DataEffect(e))
}
