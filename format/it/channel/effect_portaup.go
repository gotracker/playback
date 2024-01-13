package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PortaUp defines a portamento up effect
type PortaUp[TPeriod period.Period] DataEffect // 'F'

func (e PortaUp[TPeriod]) String() string {
	return fmt.Sprintf("F%0.2x", DataEffect(e))
}

func (e PortaUp[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	xx := mem.PortaUp(DataEffect(e))

	return m.DoChannelPortaUp(ch, period.Delta(xx)*4)
}

func (e PortaUp[TPeriod]) TraceData() string {
	return e.String()
}
