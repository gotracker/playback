package channel

import (
	"fmt"

	"github.com/gotracker/playback/mixing/sampling"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SampleOffset defines a sample offset effect
type SampleOffset[TPeriod period.Period] DataEffect // 'O'

func (e SampleOffset[TPeriod]) String() string {
	return fmt.Sprintf("O%0.2x", DataEffect(e))
}

func (e SampleOffset[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	xx := mem.SampleOffset(DataEffect(e))

	if tick != 0 {
		return nil
	}

	pos := sampling.Pos{Pos: mem.HighOffset + int(xx)*0x100}
	if mem.Shared.OldEffectMode {
		inst, err := m.GetChannelInstrument(ch)
		if err != nil {
			return err
		}
		if inst == nil || pos.Pos >= inst.GetLength().Pos {
			return nil
		}
	}
	return m.SetChannelPos(ch, pos)
}

func (e SampleOffset[TPeriod]) TraceData() string {
	return e.String()
}
