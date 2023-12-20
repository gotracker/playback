package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	"github.com/gotracker/playback/period"
)

// SetCoarsePanPosition defines a set pan position effect
type SetCoarsePanPosition[TPeriod period.Period] channel.DataEffect // 'E8x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetCoarsePanPosition[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	xy := channel.DataEffect(e)
	y := xy & 0x0F

	cs.SetPan(xmPanning.PanningFromXm(uint8(y) << 4))
	return nil
}

func (e SetCoarsePanPosition[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
