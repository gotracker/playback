package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PanningEnvelopeOff defines a panning envelope: off effect
type PanningEnvelopeOff[TPeriod period.Period] DataEffect // 'S79'

func (e PanningEnvelopeOff[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e PanningEnvelopeOff[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.SetChannelPanningEnvelopeEnable(ch, false)
}

func (e PanningEnvelopeOff[TPeriod]) TraceData() string {
	return e.String()
}
