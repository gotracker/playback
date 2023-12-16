package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
)

// SetCoarsePanPosition defines a set pan position effect
type SetCoarsePanPosition channel.DataEffect // 'E8x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetCoarsePanPosition) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	xy := channel.DataEffect(e)
	y := xy & 0x0F

	cs.SetPan(xmPanning.PanningFromXm(uint8(y) << 4))
	return nil
}

func (e SetCoarsePanPosition) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
