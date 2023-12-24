package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	"github.com/gotracker/playback/period"
)

// SetCoarsePanPosition defines a set pan position effect
type SetCoarsePanPosition[TPeriod period.Period] DataEffect // 'E8x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetCoarsePanPosition[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	xy := DataEffect(e)
	y := xy & 0x0F

	cs.GetActiveState().Pan = xmPanning.PanningFromXm(uint8(y) << 4)
	return nil
}

func (e SetCoarsePanPosition[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}
