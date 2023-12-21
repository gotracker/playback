package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// FinePortaDown defines an fine portamento down effect
type FinePortaDown ChannelCommand // 'EFx'

// Start triggers on the first tick, but before the Tick() function is called
func (e FinePortaDown) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()

	y := DataEffect(e) & 0x0F

	return doPortaDown(cs, float32(y), 4)
}

func (e FinePortaDown) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}
