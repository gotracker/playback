package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	effectIntf "github.com/gotracker/playback/format/it/effect/intf"
	"github.com/gotracker/playback/period"
)

// SetTempo defines a set tempo effect
type SetTempo[TPeriod period.Period] channel.DataEffect // 'T'

// PreStart triggers when the effect enters onto the channel state
func (e SetTempo[TPeriod]) PreStart(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	if e > 0x20 {
		m := p.(effectIntf.IT)
		if err := m.SetTempo(int(e)); err != nil {
			return err
		}
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetTempo[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e SetTempo[TPeriod]) Tick(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback, currentTick int) error {
	m := p.(effectIntf.IT)
	switch channel.DataEffect(e >> 4) {
	case 0: // decrease tempo
		if currentTick != 0 {
			mem := cs.GetMemory()
			val := int(mem.TempoDecrease(channel.DataEffect(e & 0x0F)))
			if err := m.DecreaseTempo(val); err != nil {
				return err
			}
		}
	case 1: // increase tempo
		if currentTick != 0 {
			mem := cs.GetMemory()
			val := int(mem.TempoIncrease(channel.DataEffect(e & 0x0F)))
			if err := m.IncreaseTempo(val); err != nil {
				return err
			}
		}
	default:
		if err := m.SetTempo(int(e)); err != nil {
			return err
		}
	}
	return nil
}

func (e SetTempo[TPeriod]) String() string {
	return fmt.Sprintf("T%0.2x", channel.DataEffect(e))
}
