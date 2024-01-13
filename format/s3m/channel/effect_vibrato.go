package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// Vibrato defines a vibrato effect
type Vibrato ChannelCommand // 'H'

// Start triggers on the first tick, but before the Tick() function is called
func (e Vibrato) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	return nil
}

func (e Vibrato) String() string {
	return fmt.Sprintf("H%0.2x", DataEffect(e))
}

func (e Vibrato) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	x, y := mem.Vibrato(DataEffect(e))
	// NOTE: JBC - S3M does not update on tick 0, but MOD does.
	if tick != 0 || mem.Shared.ModCompatibility {
		return withOscillatorDo(ch, m, int(x), float32(y)*4, machine.OscillatorVibrato, func(value float32) error {
			return m.SetChannelPeriodDelta(ch, period.Delta(value))
		})
	}
	return nil
}

func (e Vibrato) TraceData() string {
	return e.String()
}
