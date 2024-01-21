package channel

import (
	"fmt"

	"github.com/gotracker/gomixing/sampling"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SampleOffset defines a sample offset effect
type SampleOffset ChannelCommand // 'O'

func (e SampleOffset) String() string {
	return fmt.Sprintf("O%0.2x", DataEffect(e))
}

func (e SampleOffset) RowStart(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	xx := mem.SampleOffset(DataEffect(e))
	return m.SetChannelPos(ch, sampling.Pos{Pos: int(xx) * 0x100})
}

func (e SampleOffset) TraceData() string {
	return e.String()
}
