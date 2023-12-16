package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
	"github.com/gotracker/playback/index"
)

// OrderJump defines an order jump effect
type OrderJump ChannelCommand // 'B'

// Start triggers on the first tick, but before the Tick() function is called
func (e OrderJump) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Stop is called on the last tick of the row, but after the Tick() function is called
func (e OrderJump) Stop(cs playback.Channel[channel.Memory], p playback.Playback, lastTick int) error {
	return p.SetNextOrder(index.Order(e))
}

func (e OrderJump) String() string {
	return fmt.Sprintf("B%0.2x", channel.DataEffect(e))
}
