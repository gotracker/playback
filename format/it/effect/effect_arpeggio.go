package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// Arpeggio defines an arpeggio effect
type Arpeggio[TPeriod period.Period] channel.DataEffect // 'J'

// Start triggers on the first tick, but before the Tick() function is called
func (e Arpeggio[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	cs.SetPos(cs.GetTargetPos())
	return nil
}

// Tick is called on every tick
func (e Arpeggio[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.Arpeggio(channel.DataEffect(e))
	return doArpeggio(cs, currentTick, int8(x), int8(y))
}

func (e Arpeggio[TPeriod]) String() string {
	return fmt.Sprintf("J%0.2x", channel.DataEffect(e))
}
