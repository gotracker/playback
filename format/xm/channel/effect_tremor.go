package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// Tremor defines a tremor effect
type Tremor[TPeriod period.Period] DataEffect // 'T'

// Start triggers on the first tick, but before the Tick() function is called
func (e Tremor[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e Tremor[TPeriod]) Tick(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback, currentTick int) error {
	if currentTick != 0 {
		mem := cs.GetMemory()
		x, y := mem.Tremor(DataEffect(e))
		return doTremor(cs, currentTick, int(x)+1, int(y)+1)
	}
	return nil
}

func (e Tremor[TPeriod]) String() string {
	return fmt.Sprintf("T%0.2x", DataEffect(e))
}
