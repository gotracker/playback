package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PanningEnvelopeOn defines a panning envelope: on effect
type PanningEnvelopeOn[TPeriod period.Period] DataEffect // 'S7A'

func (e PanningEnvelopeOn[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e PanningEnvelopeOn[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.SetChannelPanningEnvelopeEnable(ch, true)
}

func (e PanningEnvelopeOn[TPeriod]) TraceData() string {
	return e.String()
}
