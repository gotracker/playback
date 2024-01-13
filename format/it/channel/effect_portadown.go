package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PortaDown defines a portamento down effect
type PortaDown[TPeriod period.Period] DataEffect // 'E'

func (e PortaDown[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e PortaDown[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	xx := mem.PortaDown(DataEffect(e))

	return m.DoChannelPortaDown(ch, period.Delta(xx)*4)
}

func (e PortaDown[TPeriod]) TraceData() string {
	return e.String()
}
