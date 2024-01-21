package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FinePortaUp defines an fine portamento up effect
type FinePortaUp[TPeriod period.Period] DataEffect // 'E1x'

func (e FinePortaUp[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e FinePortaUp[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	y := mem.FinePortaUp(DataEffect(e)) & 0x0F

	if tick != 0 {
		return nil
	}

	return m.DoChannelPortaUp(ch, period.Delta(y)*4)
}

func (e FinePortaUp[TPeriod]) TraceData() string {
	return e.String()
}
