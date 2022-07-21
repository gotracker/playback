package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
	effectIntf "github.com/gotracker/playback/format/s3m/effect/intf"
)

// SetSpeed defines a set speed effect
type SetSpeed ChannelCommand // 'A'

// PreStart triggers when the effect enters onto the channel state
func (e SetSpeed) PreStart(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	if e != 0 {
		m := p.(effectIntf.S3M)
		if err := m.SetTicks(int(e)); err != nil {
			return err
		}
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetSpeed) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e SetSpeed) String() string {
	return fmt.Sprintf("A%0.2x", channel.DataEffect(e))
}
