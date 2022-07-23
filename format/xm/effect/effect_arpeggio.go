package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// Arpeggio defines an arpeggio effect
type Arpeggio channel.DataEffect // '0'

// Start triggers on the first tick, but before the Tick() function is called
func (e Arpeggio) Start(cs *channel.State, p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	cs.SetPos(cs.GetTargetPos())
	return nil
}

// Tick is called on every tick
func (e Arpeggio) Tick(cs *channel.State, p playback.Playback, currentTick int) error {
	xy := channel.DataEffect(e)
	if xy == 0 {
		return nil
	}

	x := int8(xy >> 4)
	y := int8(xy & 0x0f)
	return doArpeggio(cs, currentTick, x, y)
}

func (e Arpeggio) String() string {
	return fmt.Sprintf("0%0.2x", channel.DataEffect(e))
}
