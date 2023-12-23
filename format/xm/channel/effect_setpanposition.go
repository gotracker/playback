package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	"github.com/gotracker/playback/period"
)

// SetPanPosition defines a set pan position effect
type SetPanPosition[TPeriod period.Period] DataEffect // '8xx'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetPanPosition[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	xx := uint8(e)

	cs.SetPan(xmPanning.PanningFromXm(xx))
	return nil
}

func (e SetPanPosition[TPeriod]) String() string {
	return fmt.Sprintf("8%0.2x", DataEffect(e))
}
