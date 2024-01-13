package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SurroundOn defines a set surround on effect
type SurroundOn ChannelCommand // 'S91'

func (e SurroundOn) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e SurroundOn) RowStart(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	// TODO: support for center rear panning
	return nil
}

func (e SurroundOn) TraceData() string {
	return e.String()
}
