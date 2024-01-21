package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// NoteCut defines a note cut effect
type NoteCut ChannelCommand // 'SCx'

func (e NoteCut) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e NoteCut) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	if tick == 0 {
		return nil
	}

	x := DataEffect(e) & 0xf

	if tick == int(x) {
		return m.ChannelStop(ch)
	}
	return nil
}

func (e NoteCut) TraceData() string {
	return e.String()
}
