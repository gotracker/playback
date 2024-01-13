package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetSpeed defines a set speed effect
type SetSpeed ChannelCommand // 'A'

func (e SetSpeed) String() string {
	return fmt.Sprintf("A%0.2x", DataEffect(e))
}

func (e SetSpeed) RowStart(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	return m.SetTempo(int(e))
}

func (e SetSpeed) TraceData() string {
	return e.String()
}
