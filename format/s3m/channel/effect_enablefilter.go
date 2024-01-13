package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// EnableFilter defines a set filter enable effect
type EnableFilter ChannelCommand // 'S0x'

func (e EnableFilter) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e EnableFilter) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	if tick != 0 {
		return nil
	}

	x := DataEffect(e) & 0xf
	return m.SetFilterOnAllChannelsByFilterName("amigalpf", x != 0)
}

func (e EnableFilter) TraceData() string {
	return e.String()
}
