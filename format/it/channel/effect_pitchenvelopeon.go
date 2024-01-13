package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PitchEnvelopeOn defines a panning envelope: on effect
type PitchEnvelopeOn[TPeriod period.Period] DataEffect // 'S7C'

func (e PitchEnvelopeOn[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e PitchEnvelopeOn[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.SetChannelPitchEnvelopeEnable(ch, true)
}

func (e PitchEnvelopeOn[TPeriod]) TraceData() string {
	return e.String()
}
