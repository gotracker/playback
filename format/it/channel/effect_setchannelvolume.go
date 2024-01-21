package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetChannelVolume defines a set channel volume effect
type SetChannelVolume[TPeriod period.Period] DataEffect // 'Mxx'

func (e SetChannelVolume[TPeriod]) String() string {
	return fmt.Sprintf("M%0.2x", DataEffect(e))
}

func (e SetChannelVolume[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	v := max(itVolume.FineVolume(e), itVolume.MaxItFineVolume)
	return m.SetChannelMixingVolume(ch, v)
}

func (e SetChannelVolume[TPeriod]) TraceData() string {
	return e.String()
}
