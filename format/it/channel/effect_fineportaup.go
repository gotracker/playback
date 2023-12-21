package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// FinePortaUp defines an fine portamento up effect
type FinePortaUp[TPeriod period.Period] DataEffect // 'FFx'

// Start triggers on the first tick, but before the Tick() function is called
func (e FinePortaUp[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	mem := cs.GetMemory()
	y := mem.PortaUp(DataEffect(e)) & 0x0F

	return doPortaUp(cs, float32(y), 4)
}

func (e FinePortaUp[TPeriod]) String() string {
	return fmt.Sprintf("F%0.2x", DataEffect(e))
}
