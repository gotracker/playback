package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FinePortaDown defines an fine portamento down effect
type FinePortaDown ChannelCommand // 'EFx'

func (e FinePortaDown) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e FinePortaDown) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	if tick != 0 {
		return nil
	}

	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	y := DataEffect(e) & 0x0F

	var mul period.Delta = 4
	if mem.Shared.ST300Portas {
		mul = 2
	}
	return m.DoChannelPortaDown(ch, period.Delta(y)*mul)
}

func (e FinePortaDown) TraceData() string {
	return e.String()
}
