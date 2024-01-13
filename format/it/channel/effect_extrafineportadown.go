package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// ExtraFinePortaDown defines an extra-fine portamento down effect
type ExtraFinePortaDown[TPeriod period.Period] DataEffect // 'EEx'

func (e ExtraFinePortaDown[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e ExtraFinePortaDown[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	y := mem.PortaDown(DataEffect(e)) & 0x0F

	if tick != 0 {
		return nil
	}

	return m.DoChannelPortaDown(ch, period.Delta(y)*1)
}

func (e ExtraFinePortaDown[TPeriod]) TraceData() string {
	return e.String()
}
