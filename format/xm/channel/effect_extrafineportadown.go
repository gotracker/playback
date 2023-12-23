package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// ExtraFinePortaDown defines an extra-fine portamento down effect
type ExtraFinePortaDown[TPeriod period.Period] DataEffect // 'X2x'

// Start triggers on the first tick, but before the Tick() function is called
func (e ExtraFinePortaDown[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	mem := cs.GetMemory()
	xx := mem.ExtraFinePortaDown(DataEffect(e))
	y := xx & 0x0F

	return doPortaDown(cs, float32(y), 1)
}

func (e ExtraFinePortaDown[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}
