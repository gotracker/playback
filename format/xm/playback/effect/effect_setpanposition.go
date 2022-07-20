package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	xmPanning "github.com/gotracker/playback/format/xm/conversion/panning"
	"github.com/gotracker/playback/format/xm/layout/channel"
)

// SetPanPosition defines a set pan position effect
type SetPanPosition channel.DataEffect // '8xx'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetPanPosition) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	xx := uint8(e)

	cs.SetPan(xmPanning.PanningFromXm(xx))
	return nil
}

func (e SetPanPosition) String() string {
	return fmt.Sprintf("8%0.2x", channel.DataEffect(e))
}
