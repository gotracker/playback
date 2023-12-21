package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// PatternDelay defines a pattern delay effect
type PatternDelay ChannelCommand // 'SEx'

// PreStart triggers when the effect enters onto the channel state
func (e PatternDelay) PreStart(cs S3MChannel, p playback.Playback) error {
	m := p.(S3M)
	return m.SetPatternDelay(int(DataEffect(e) & 0x0F))
}

// Start triggers on the first tick, but before the Tick() function is called
func (e PatternDelay) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e PatternDelay) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
