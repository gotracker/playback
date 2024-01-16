package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PortaToNote defines a portamento-to-note effect
type PortaToNote[TPeriod period.Period] DataEffect // 'G'

func (e PortaToNote[TPeriod]) String() string {
	return fmt.Sprintf("G%0.2x", DataEffect(e))
}

func (e PortaToNote[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.StartChannelPortaToNote(ch)
}

func (e PortaToNote[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	xx := mem.PortaToNote(DataEffect(e))

	if !mem.Shared.OldEffectMode || tick != 0 {
		return m.DoChannelPortaToNote(ch, period.Delta(xx)*4, false)
	}
	return nil
}

func (e PortaToNote[TPeriod]) TraceData() string {
	return e.String()
}
