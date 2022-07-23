package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	effectIntf "github.com/gotracker/playback/format/xm/effect/intf"
)

// GlobalVolumeSlide defines a global volume slide effect
type GlobalVolumeSlide channel.DataEffect // 'H'

// Start triggers on the first tick, but before the Tick() function is called
func (e GlobalVolumeSlide) Start(cs *channel.State, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e GlobalVolumeSlide) Tick(cs *channel.State, p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.GlobalVolumeSlide(channel.DataEffect(e))

	if currentTick == 0 {
		return nil
	}

	m := p.(effectIntf.XM)

	if x == 0 {
		// global vol slide down
		return doGlobalVolSlide(m, -float32(y), 1.0)
	} else if y == 0 {
		// global vol slide up
		return doGlobalVolSlide(m, float32(y), 1.0)
	}
	return nil
}

func (e GlobalVolumeSlide) String() string {
	return fmt.Sprintf("H%0.2x", channel.DataEffect(e))
}
