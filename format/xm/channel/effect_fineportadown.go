package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// FinePortaDown defines an fine portamento down effect
type FinePortaDown[TPeriod period.Period] DataEffect // 'E2x'

// Start triggers on the first tick, but before the Tick() function is called
func (e FinePortaDown[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	mem := cs.GetMemory()
	xy := mem.FinePortaDown(DataEffect(e))
	y := xy & 0x0F

	return doPortaDown(cs, float32(y), 4)
}

func (e FinePortaDown[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}
