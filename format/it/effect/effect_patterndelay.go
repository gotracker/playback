package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	effectIntf "github.com/gotracker/playback/format/it/effect/intf"
	"github.com/gotracker/playback/period"
)

// PatternDelay defines a pattern delay effect
type PatternDelay[TPeriod period.Period] channel.DataEffect // 'SEx'

// PreStart triggers when the effect enters onto the channel state
func (e PatternDelay[TPeriod]) PreStart(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	m := p.(effectIntf.IT)
	return m.SetPatternDelay(int(channel.DataEffect(e) & 0x0F))
}

// Start triggers on the first tick, but before the Tick() function is called
func (e PatternDelay[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e PatternDelay[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
