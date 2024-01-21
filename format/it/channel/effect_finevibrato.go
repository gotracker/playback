package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FineVibrato defines an fine vibrato effect
type FineVibrato[TPeriod period.Period] DataEffect // 'U'

func (e FineVibrato[TPeriod]) String() string {
	return fmt.Sprintf("U%0.2x", DataEffect(e))
}

func (e FineVibrato[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	x, y := mem.Vibrato(DataEffect(e))

	if tick == 0 {
		return nil
	}

	return withOscillatorDo(ch, m, int(x), float32(y)*1, machine.OscillatorVibrato, func(value float32) error {
		return m.SetChannelPeriodDelta(ch, period.Delta(value))
	})
}

func (e FineVibrato[TPeriod]) TraceData() string {
	return e.String()
}
