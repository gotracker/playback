package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// VolumeEnvelopeOff defines a volume envelope: off effect
type VolumeEnvelopeOff[TPeriod period.Period] DataEffect // 'S77'

func (e VolumeEnvelopeOff[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e VolumeEnvelopeOff[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.SetChannelVolumeEnvelopeEnable(ch, false)
}

func (e VolumeEnvelopeOff[TPeriod]) TraceData() string {
	return e.String()
}
