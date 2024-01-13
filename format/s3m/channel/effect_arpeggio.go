package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// Arpeggio defines an arpeggio effect
type Arpeggio ChannelCommand // 'J'

func (e Arpeggio) String() string {
	return fmt.Sprintf("J%0.2x", DataEffect(e))
}

func (e Arpeggio) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x, y := mem.LastNonZeroXY(DataEffect(e))
	return doArpeggio(ch, m, tick, int8(x), int8(y))
}

func (e Arpeggio) TraceData() string {
	return e.String()
}
