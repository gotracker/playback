package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	"github.com/gotracker/playback/format/it/system"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FinePortaDown defines an fine portamento down effect
type FinePortaDown[TPeriod period.Period] DataEffect // 'EFx'

func (e FinePortaDown[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e FinePortaDown[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	if tick != 0 {
		return nil
	}

	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	y := mem.PortaDown(DataEffect(e)) & 0x0F

	return m.DoChannelPortaDown(ch, period.Delta(y)*system.SlideFinesPerSemitone)
}

func (e FinePortaDown[TPeriod]) TraceData() string {
	return e.String()
}
