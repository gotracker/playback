package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FinePortaDown defines an fine portamento down effect
type FinePortaDown[TPeriod period.Period] DataEffect // 'E2x'

func (e FinePortaDown[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e FinePortaDown[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	y := mem.FinePortaDown(DataEffect(e)) & 0x0F

	if tick != 0 {
		return nil
	}

	return m.DoChannelPortaDown(ch, period.Delta(y)*4)
}

func (e FinePortaDown[TPeriod]) TraceData() string {
	return e.String()
}
