package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// FinePortaUp defines an fine portamento up effect
type FinePortaUp ChannelCommand // 'FFx'

// Start triggers on the first tick, but before the Tick() function is called
func (e FinePortaUp) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	y := DataEffect(e) & 0x0F

	return doPortaUp(cs, float32(y), 4)
}

func (e FinePortaUp) String() string {
	return fmt.Sprintf("F%0.2x", DataEffect(e))
}
