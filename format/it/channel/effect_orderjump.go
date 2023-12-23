package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
)

// OrderJump defines an order jump effect
type OrderJump[TPeriod period.Period] DataEffect // 'B'

// Start triggers on the first tick, but before the Tick() function is called
func (e OrderJump[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Stop is called on the last tick of the row, but after the Tick() function is called
func (e OrderJump[TPeriod]) Stop(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback, lastTick int) error {
	return p.SetNextOrder(index.Order(e))
}

func (e OrderJump[TPeriod]) String() string {
	return fmt.Sprintf("B%0.2x", DataEffect(e))
}
