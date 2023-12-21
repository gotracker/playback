package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// HighOffset defines a sample high offset effect
type HighOffset[TPeriod period.Period] DataEffect // 'SAx'

// Start triggers on the first tick, but before the Tick() function is called
func (e HighOffset[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	mem := cs.GetMemory()

	xx := DataEffect(e)

	mem.HighOffset = int(xx) * 0x10000
	return nil
}

func (e HighOffset[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
