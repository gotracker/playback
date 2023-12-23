package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// FinePortaUp defines an fine portamento up effect
type FinePortaUp[TPeriod period.Period] DataEffect // 'E1x'

// Start triggers on the first tick, but before the Tick() function is called
func (e FinePortaUp[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	mem := cs.GetMemory()
	xy := mem.FinePortaUp(DataEffect(e))
	y := xy & 0x0F

	return doPortaUp(cs, float32(y), 4)
}

func (e FinePortaUp[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}
