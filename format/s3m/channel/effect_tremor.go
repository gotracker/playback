package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// Tremor defines a tremor effect
type Tremor ChannelCommand // 'I'

func (e Tremor) String() string {
	return fmt.Sprintf("I%0.2x", DataEffect(e))
}

func (e Tremor) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	x, y := mem.LastNonZeroXY(DataEffect(e))
	return doTremor(ch, m, int(x)+1, int(y)+1)
}

func (e Tremor) TraceData() string {
	return e.String()
}
