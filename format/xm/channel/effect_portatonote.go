package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PortaToNote defines a portamento-to-note effect
type PortaToNote[TPeriod period.Period] DataEffect // '3'

func (e PortaToNote[TPeriod]) String() string {
	return fmt.Sprintf("3%0.2x", DataEffect(e))
}

func (e PortaToNote[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	return m.StartChannelPortaToNote(ch)
}

func (e PortaToNote[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	if tick == 0 {
		return nil
	}

	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	xx := mem.PortaToNote(DataEffect(e))
	return m.DoChannelPortaToNote(ch, period.Delta(xx)*4, false)
}

func (e PortaToNote[TPeriod]) TraceData() string {
	return e.String()
}
