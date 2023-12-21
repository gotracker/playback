package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// SetTempo defines a set tempo effect
type SetTempo[TPeriod period.Period] DataEffect // 'F'

// PreStart triggers when the effect enters onto the channel state
func (e SetTempo[TPeriod]) PreStart(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	if e > 0x20 {
		m := p.(XM)
		if err := m.SetTempo(int(e)); err != nil {
			return err
		}
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetTempo[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e SetTempo[TPeriod]) Tick(cs playback.Channel[TPeriod, Memory], p playback.Playback, currentTick int) error {
	m := p.(XM)
	return m.SetTempo(int(e))
}

func (e SetTempo[TPeriod]) String() string {
	return fmt.Sprintf("F%0.2x", DataEffect(e))
}
