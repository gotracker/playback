package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// FinePortaDown defines an fine portamento down effect
type FinePortaDown channel.DataEffect // 'E2x'

// Start triggers on the first tick, but before the Tick() function is called
func (e FinePortaDown) Start(cs *channel.State, p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	mem := cs.GetMemory()
	xy := mem.FinePortaDown(channel.DataEffect(e))
	y := xy & 0x0F

	return doPortaDown(cs, float32(y), 4)
}

func (e FinePortaDown) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
