package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// FinePortaUp defines an fine portamento up effect
type FinePortaUp[TPeriod period.Period] channel.DataEffect // 'FFx'

// Start triggers on the first tick, but before the Tick() function is called
func (e FinePortaUp[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	mem := cs.GetMemory()
	y := mem.PortaUp(channel.DataEffect(e)) & 0x0F

	return doPortaUp(cs, float32(y), 4)
}

func (e FinePortaUp[TPeriod]) String() string {
	return fmt.Sprintf("F%0.2x", channel.DataEffect(e))
}
