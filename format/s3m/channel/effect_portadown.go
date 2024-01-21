package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PortaDown defines a portamento down effect
type PortaDown ChannelCommand // 'E'

func (e PortaDown) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e PortaDown) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	xx := mem.Porta(DataEffect(e))

	if tick == 0 && !mem.Shared.AmigaSlides {
		return nil
	}

	var mul period.Delta = 4
	if mem.Shared.ST300Portas {
		mul = 2
	}
	return m.DoChannelPortaDown(ch, period.Delta(xx)*mul)
}

func (e PortaDown) TraceData() string {
	return e.String()
}
