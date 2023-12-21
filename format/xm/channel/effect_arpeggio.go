package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// Arpeggio defines an arpeggio effect
type Arpeggio[TPeriod period.Period] DataEffect // '0'

// Start triggers on the first tick, but before the Tick() function is called
func (e Arpeggio[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	cs.SetPos(cs.GetTargetPos())
	return nil
}

// Tick is called on every tick
func (e Arpeggio[TPeriod]) Tick(cs playback.Channel[TPeriod, Memory], p playback.Playback, currentTick int) error {
	xy := DataEffect(e)
	if xy == 0 {
		return nil
	}

	x := int8(xy >> 4)
	y := int8(xy & 0x0f)
	return doArpeggio(cs, currentTick, x, y)
}

func (e Arpeggio[TPeriod]) String() string {
	return fmt.Sprintf("0%0.2x", DataEffect(e))
}
