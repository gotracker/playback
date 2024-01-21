package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetPanPosition defines a set pan position effect
type SetPanPosition[TPeriod period.Period] DataEffect // 'Xxx'

func (e SetPanPosition[TPeriod]) String() string {
	return fmt.Sprintf("X%0.2x", DataEffect(e))
}

func (e SetPanPosition[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	pan := itPanning.Panning(e)
	return m.SetChannelPan(ch, pan)
}

func (e SetPanPosition[TPeriod]) TraceData() string {
	return e.String()
}
