package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
	effectIntf "github.com/gotracker/playback/format/s3m/playback/effect/intf"
)

// EnableFilter defines a set filter enable effect
type EnableFilter ChannelCommand // 'S0x'

// Start triggers on the first tick, but before the Tick() function is called
func (e EnableFilter) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := channel.DataEffect(e) & 0xf
	on := x != 0

	pb := p.(effectIntf.S3M)
	pb.SetFilterEnable(on)
	return nil
}

func (e EnableFilter) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
