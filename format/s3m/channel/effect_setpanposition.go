package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetPanPosition defines a set pan position effect
type SetPanPosition ChannelCommand // 'S8x'

func (e SetPanPosition) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e SetPanPosition) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	return m.SetChannelPan(ch, s3mPanning.Panning(uint8(e)&0xf))
}

func (e SetPanPosition) TraceData() string {
	return e.String()
}
