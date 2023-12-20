package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	effectIntf "github.com/gotracker/playback/format/it/effect/intf"
	"github.com/gotracker/playback/period"
)

// SetSpeed defines a set speed effect
type SetSpeed[TPeriod period.Period] channel.DataEffect // 'A'

// PreStart triggers when the effect enters onto the channel state
func (e SetSpeed[TPeriod]) PreStart(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	if e != 0 {
		m := p.(effectIntf.IT)
		if err := m.SetTicks(int(e)); err != nil {
			return err
		}
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetSpeed[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e SetSpeed[TPeriod]) String() string {
	return fmt.Sprintf("A%0.2x", channel.DataEffect(e))
}
