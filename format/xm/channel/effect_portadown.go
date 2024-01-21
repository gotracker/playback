package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PortaDown defines a portamento down effect
type PortaDown[TPeriod period.Period] DataEffect // '2'

func (e PortaDown[TPeriod]) String() string {
	return fmt.Sprintf("2%0.2x", DataEffect(e))
}

func (e PortaDown[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	xx := mem.PortaDown(DataEffect(e))

	if tick == 0 {
		return nil
	}

	return m.DoChannelPortaDown(ch, period.Delta(xx)*4)
}

func (e PortaDown[TPeriod]) TraceData() string {
	return e.String()
}
