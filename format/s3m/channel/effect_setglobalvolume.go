package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetGlobalVolume defines a set global volume effect
type SetGlobalVolume ChannelCommand // 'V'

func (e SetGlobalVolume) String() string {
	return fmt.Sprintf("V%0.2x", DataEffect(e))
}

func (e SetGlobalVolume) RowStart(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	return m.SetGlobalVolume(s3mVolume.Volume(DataEffect(e)))
}

func (e SetGlobalVolume) TraceData() string {
	return e.String()
}
