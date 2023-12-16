package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	effectIntf "github.com/gotracker/playback/format/xm/effect/intf"
)

// SetSpeed defines a set speed effect
type SetSpeed channel.DataEffect // 'F'

// PreStart triggers when the effect enters onto the channel state
func (e SetSpeed) PreStart(cs playback.Channel[channel.Memory], p playback.Playback) error {
	if e != 0 {
		m := p.(effectIntf.XM)
		if err := m.SetTicks(int(e)); err != nil {
			return err
		}
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetSpeed) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e SetSpeed) String() string {
	return fmt.Sprintf("F%0.2x", channel.DataEffect(e))
}
