package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
)

// SetPanPosition defines a set pan position effect
type SetPanPosition channel.DataEffect // '8xx'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetPanPosition) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	xx := uint8(e)

	cs.SetPan(xmPanning.PanningFromXm(xx))
	return nil
}

func (e SetPanPosition) String() string {
	return fmt.Sprintf("8%0.2x", channel.DataEffect(e))
}
