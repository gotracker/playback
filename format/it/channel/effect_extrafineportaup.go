package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// ExtraFinePortaUp defines an extra-fine portamento up effect
type ExtraFinePortaUp[TPeriod period.Period] DataEffect // 'FEx'

func (e ExtraFinePortaUp[TPeriod]) String() string {
	return fmt.Sprintf("F%0.2x", DataEffect(e))
}

func (e ExtraFinePortaUp[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	y := mem.PortaUp(DataEffect(e)) & 0x0F

	if tick != 0 {
		return nil
	}
	return m.DoChannelPortaUp(ch, period.Delta(y)*1)
}

func (e ExtraFinePortaUp[TPeriod]) TraceData() string {
	return e.String()
}
