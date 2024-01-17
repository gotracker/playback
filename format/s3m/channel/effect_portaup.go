package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PortaUp defines a portamento up effect
type PortaUp ChannelCommand // 'F'

func (e PortaUp) String() string {
	return fmt.Sprintf("F%0.2x", DataEffect(e))
}

func (e PortaUp) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	xx := mem.Porta(DataEffect(e))

	if tick == 0 && !mem.Shared.AmigaSlides {
		return nil
	}

	return m.DoChannelPortaUp(ch, period.Delta(xx)*4)
}

func (e PortaUp) TraceData() string {
	return e.String()
}
