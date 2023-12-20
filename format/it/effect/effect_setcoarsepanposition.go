package effect

import (
	"fmt"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	itPanning "github.com/gotracker/playback/format/it/panning"
	"github.com/gotracker/playback/period"
)

// SetCoarsePanPosition defines a set coarse pan position effect
type SetCoarsePanPosition[TPeriod period.Period] channel.DataEffect // 'S8x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetCoarsePanPosition[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := channel.DataEffect(e) & 0xf

	pan := itfile.PanValue(x << 2)

	cs.SetPan(itPanning.FromItPanning(pan))
	return nil
}

func (e SetCoarsePanPosition[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
