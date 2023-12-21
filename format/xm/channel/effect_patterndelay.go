package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// PatternDelay defines a pattern delay effect
type PatternDelay[TPeriod period.Period] DataEffect // 'SEx'

// PreStart triggers when the effect enters onto the channel state
func (e PatternDelay[TPeriod]) PreStart(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	m := p.(XM)
	return m.SetPatternDelay(int(DataEffect(e) & 0x0F))
}

// Start triggers on the first tick, but before the Tick() function is called
func (e PatternDelay[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e PatternDelay[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
