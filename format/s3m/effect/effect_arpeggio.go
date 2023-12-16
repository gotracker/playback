package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
)

// Arpeggio defines an arpeggio effect
type Arpeggio ChannelCommand // 'J'

// Start triggers on the first tick, but before the Tick() function is called
func (e Arpeggio) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	cs.SetPos(cs.GetTargetPos())
	return nil
}

// Tick is called on every tick
func (e Arpeggio) Tick(cs playback.Channel[channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.LastNonZeroXY(channel.DataEffect(e))
	return doArpeggio(cs, currentTick, int8(x), int8(y))
}

func (e Arpeggio) String() string {
	return fmt.Sprintf("J%0.2x", channel.DataEffect(e))
}
