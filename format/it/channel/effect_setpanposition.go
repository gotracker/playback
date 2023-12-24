package channel

import (
	"fmt"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"

	"github.com/gotracker/playback"
	itPanning "github.com/gotracker/playback/format/it/panning"
	"github.com/gotracker/playback/period"
)

// SetPanPosition defines a set pan position effect
type SetPanPosition[TPeriod period.Period] DataEffect // 'Xxx'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetPanPosition[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := DataEffect(e)

	pan := itfile.PanValue(x)

	cs.GetActiveState().Pan = itPanning.FromItPanning(pan)
	return nil
}

func (e SetPanPosition[TPeriod]) String() string {
	return fmt.Sprintf("X%0.2x", DataEffect(e))
}
