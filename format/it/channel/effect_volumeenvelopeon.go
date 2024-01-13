package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// VolumeEnvelopeOn defines a volume envelope: on effect
type VolumeEnvelopeOn[TPeriod period.Period] DataEffect // 'S78'

func (e VolumeEnvelopeOn[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e VolumeEnvelopeOn[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.SetChannelVolumeEnvelopeEnable(ch, true)
}

func (e VolumeEnvelopeOn[TPeriod]) TraceData() string {
	return e.String()
}
