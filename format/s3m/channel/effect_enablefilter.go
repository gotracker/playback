package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// EnableFilter defines a set filter enable effect
type EnableFilter ChannelCommand // 'S0x'

// Start triggers on the first tick, but before the Tick() function is called
func (e EnableFilter) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := DataEffect(e) & 0xf
	on := x != 0

	pb := p.(S3M)
	pb.SetFilterEnable(on)
	return nil
}

func (e EnableFilter) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
