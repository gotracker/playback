package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	effectIntf "github.com/gotracker/playback/format/xm/playback/effect/intf"
)

// SetTempo defines a set tempo effect
type SetTempo channel.DataEffect // 'F'

// PreStart triggers when the effect enters onto the channel state
func (e SetTempo) PreStart(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	if e > 0x20 {
		m := p.(effectIntf.XM)
		if err := m.SetTempo(int(e)); err != nil {
			return err
		}
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetTempo) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e SetTempo) Tick(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback, currentTick int) error {
	m := p.(effectIntf.XM)
	return m.SetTempo(int(e))
}

func (e SetTempo) String() string {
	return fmt.Sprintf("F%0.2x", channel.DataEffect(e))
}
