package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// StereoControl defines a set stereo control effect
type StereoControl ChannelCommand // 'SAx'

func (e StereoControl) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e StereoControl) RowStart(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	x := uint8(e) & 0xf

	if x > 7 {
		return m.SetChannelPan(ch, s3mPanning.Panning(x-8))
	} else {
		return m.SetChannelPan(ch, s3mPanning.Panning(x+8))
	}
}

func (e StereoControl) TraceData() string {
	return e.String()
}
