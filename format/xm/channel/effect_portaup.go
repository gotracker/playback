package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PortaUp defines a portamento up effect
type PortaUp[TPeriod period.Period] DataEffect // '1'

func (e PortaUp[TPeriod]) String() string {
	return fmt.Sprintf("1%0.2x", DataEffect(e))
}

func (e PortaUp[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	xx := mem.PortaUp(DataEffect(e))

	if tick == 0 {
		return nil
	}

	return m.DoChannelPortaUp(ch, period.Delta(xx)*4)
}

func (e PortaUp[TPeriod]) TraceData() string {
	return e.String()
}
