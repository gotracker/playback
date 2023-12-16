package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
)

// Tremor defines a tremor effect
type Tremor channel.DataEffect // 'I'

// Start triggers on the first tick, but before the Tick() function is called
func (e Tremor) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

// Tick is called on every tick
func (e Tremor) Tick(cs playback.Channel[channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	x, y := mem.Tremor(channel.DataEffect(e))
	return doTremor(cs, currentTick, int(x)+1, int(y)+1)
}

func (e Tremor) String() string {
	return fmt.Sprintf("I%0.2x", channel.DataEffect(e))
}
