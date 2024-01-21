package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PitchEnvelopeOff defines a panning envelope: off effect
type PitchEnvelopeOff[TPeriod period.Period] DataEffect // 'S7B'

func (e PitchEnvelopeOff[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e PitchEnvelopeOff[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.SetChannelPitchEnvelopeEnable(ch, false)
}

func (e PitchEnvelopeOff[TPeriod]) TraceData() string {
	return e.String()
}
