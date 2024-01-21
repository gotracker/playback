package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FineVibrato defines an fine vibrato effect
type FineVibrato ChannelCommand // 'U'

func (e FineVibrato) String() string {
	return fmt.Sprintf("U%0.2x", DataEffect(e))
}

func (e FineVibrato) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x, y := mem.Vibrato(DataEffect(e))
	// NOTE: JBC - S3M does not update on tick 0, but MOD does.
	if tick != 0 || mem.Shared.ModCompatibility {
		return withOscillatorDo(ch, m, int(x), float32(y)*1, machine.OscillatorVibrato, func(value float32) error {
			return m.SetChannelPeriodDelta(ch, period.Delta(value))
		})
	}
	return nil
}

func (e FineVibrato) TraceData() string {
	return e.String()
}
