package channel

import (
	"fmt"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"

	"github.com/gotracker/playback"
	itPanning "github.com/gotracker/playback/format/it/panning"
	"github.com/gotracker/playback/period"
)

// SetCoarsePanPosition defines a set coarse pan position effect
type SetCoarsePanPosition[TPeriod period.Period] DataEffect // 'S8x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetCoarsePanPosition[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := DataEffect(e) & 0xf

	pan := itfile.PanValue(x << 2)

	cs.SetPan(itPanning.FromItPanning(pan))
	return nil
}

func (e SetCoarsePanPosition[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
