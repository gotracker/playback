package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PortaToNote defines a portamento-to-note effect
type PortaToNote ChannelCommand // 'G'

func (e PortaToNote) String() string {
	return fmt.Sprintf("G%0.2x", DataEffect(e))
}

func (e PortaToNote) RowStart(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	return m.StartChannelPortaToNote(ch)
}

func (e PortaToNote) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	xx := mem.PortaToNote(DataEffect(e))

	if tick == 0 && !mem.Shared.AmigaSlides {
		return nil
	}

	return m.DoChannelPortaToNote(ch, period.Delta(xx)*4, true)
}

func (e PortaToNote) TraceData() string {
	return e.String()
}
