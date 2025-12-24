package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/voice/types"
)

// Panbrello defines a panning 'vibrato' effect
type Panbrello[TPeriod period.Period] DataEffect // 'Y'

func (e Panbrello[TPeriod]) String() string {
	return fmt.Sprintf("Y%0.2x", DataEffect(e))
}

func (e Panbrello[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x, y := mem.Panbrello(DataEffect(e))

	mul := float32(4)
	if mem.Shared.OldEffectMode {
		if tick == 0 {
			return nil
		}
		mul = 8
	}
	return withOscillatorDo(ch, m, int(x), float32(y)*mul, machine.OscillatorPanbrello, func(value float32) error {
		return m.SetChannelPanningDelta(ch, types.PanDelta(value))
	})
}

func (e Panbrello[TPeriod]) TraceData() string {
	return e.String()
}
