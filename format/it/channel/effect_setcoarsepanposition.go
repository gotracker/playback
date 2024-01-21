package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetCoarsePanPosition defines a set coarse pan position effect
type SetCoarsePanPosition[TPeriod period.Period] DataEffect // 'S8x'

func (e SetCoarsePanPosition[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e SetCoarsePanPosition[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	pan := itPanning.Panning((e & 0x0f) << 2)
	return m.SetChannelPan(ch, pan)
}

func (e SetCoarsePanPosition[TPeriod]) TraceData() string {
	return e.String()
}
