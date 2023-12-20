package effect

import (
	"fmt"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	itPanning "github.com/gotracker/playback/format/it/panning"
	"github.com/gotracker/playback/period"
)

// SetPanPosition defines a set pan position effect
type SetPanPosition[TPeriod period.Period] channel.DataEffect // 'Xxx'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetPanPosition[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := channel.DataEffect(e)

	pan := itfile.PanValue(x)

	cs.SetPan(itPanning.FromItPanning(pan))
	return nil
}

func (e SetPanPosition[TPeriod]) String() string {
	return fmt.Sprintf("X%0.2x", channel.DataEffect(e))
}
