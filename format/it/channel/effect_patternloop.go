package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PatternLoop defines a pattern loop effect
type PatternLoop[TPeriod period.Period] DataEffect // 'SBx'

func (e PatternLoop[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e PatternLoop[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	x := DataEffect(e) & 0x0F

	if x == 0 {
		// set loop start
		return m.SetPatternLoopStart(ch)
	} else {
		// set loop end + count
		return m.SetPatternLoops(ch, int(x))
	}
}

func (e PatternLoop[TPeriod]) TraceData() string {
	return e.String()
}
