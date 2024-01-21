package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PatternDelay defines a pattern delay effect
type PatternDelay ChannelCommand // 'SEx'

func (e PatternDelay) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e PatternDelay) RowStart(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	times := int(DataEffect(e) & 0x0F)
	return m.RowRepeat(times)
}

func (e PatternDelay) TraceData() string {
	return e.String()
}
