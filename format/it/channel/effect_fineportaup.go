package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FinePortaUp defines an fine portamento up effect
type FinePortaUp[TPeriod period.Period] DataEffect // 'FFx'

func (e FinePortaUp[TPeriod]) String() string {
	return fmt.Sprintf("F%0.2x", DataEffect(e))
}

func (e FinePortaUp[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	y := mem.PortaUp(DataEffect(e))

	if tick != 0 {
		return nil
	}

	return m.DoChannelPortaUp(ch, period.Delta(y)*4)
}

func (e FinePortaUp[TPeriod]) TraceData() string {
	return e.String()
}
