package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PatternLoop defines a pattern loop effect
type PatternLoop ChannelCommand // 'SBx'

func (e PatternLoop) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e PatternLoop) RowStart(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	x := DataEffect(e) & 0x0F

	if x == 0 {
		m.SetPatternLoopStart(ch)
	} else {
		m.SetPatternLoops(ch, int(x))
	}
	return nil
}

func (e PatternLoop) TraceData() string {
	return e.String()
}
