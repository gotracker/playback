package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// Vibrato defines a vibrato effect
type Vibrato[TPeriod period.Period] DataEffect // 'H'

func (e Vibrato[TPeriod]) String() string {
	return fmt.Sprintf("H%0.2x", DataEffect(e))
}

func (e Vibrato[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x, y := mem.Vibrato(DataEffect(e))

	mul := float32(4)
	if mem.Shared.OldEffectMode {
		if tick == 0 {
			return nil
		}
		mul = 8
	}
	return withOscillatorDo(ch, m, int(x), float32(y)*mul, machine.OscillatorVibrato, func(value float32) error {
		return m.SetChannelPeriodDelta(ch, period.Delta(value))
	})
}

func (e Vibrato[TPeriod]) TraceData() string {
	return e.String()
}
