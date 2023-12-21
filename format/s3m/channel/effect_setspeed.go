package channel

import (
	"fmt"

	"github.com/gotracker/playback"
)

// SetSpeed defines a set speed effect
type SetSpeed ChannelCommand // 'A'

// PreStart triggers when the effect enters onto the channel state
func (e SetSpeed) PreStart(cs S3MChannel, p playback.Playback) error {
	if e != 0 {
		m := p.(S3M)
		if err := m.SetTicks(int(e)); err != nil {
			return err
		}
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetSpeed) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e SetSpeed) String() string {
	return fmt.Sprintf("A%0.2x", DataEffect(e))
}
