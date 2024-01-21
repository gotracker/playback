package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// Arpeggio defines an arpeggio effect
type Arpeggio[TPeriod period.Period] DataEffect // 'J'

func (e Arpeggio[TPeriod]) String() string {
	return fmt.Sprintf("J%0.2x", DataEffect(e))
}

func (e Arpeggio[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	x, y := mem.Arpeggio(DataEffect(e))
	return doArpeggio(ch, m, tick, int8(x), int8(y))
}

func (e Arpeggio[TPeriod]) TraceData() string {
	return e.String()
}
