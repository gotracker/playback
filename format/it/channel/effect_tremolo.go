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

// Tremolo defines a tremolo effect
type Tremolo[TPeriod period.Period] DataEffect // 'R'

func (e Tremolo[TPeriod]) String() string {
	return fmt.Sprintf("R%0.2x", DataEffect(e))
}

func (e Tremolo[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	x, y := mem.Tremolo(DataEffect(e))

	// NOTE: JBC - IT dos not update on tick 0, but MOD does.
	// Maybe need to add a flag for converted MOD backward compatibility?
	if tick == 0 {
		return nil
	}

	return withOscillatorDo(ch, m, int(x), float32(y)*4, machine.OscillatorTremolo, func(value float32) error {
		return m.SetChannelVolumeDelta(ch, types.VolumeDelta(value))
	})
}

func (e Tremolo[TPeriod]) TraceData() string {
	return e.String()
}
