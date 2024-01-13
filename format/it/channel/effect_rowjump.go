package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// RowJump defines a row jump effect
type RowJump[TPeriod period.Period] DataEffect // 'C'

func (e RowJump[TPeriod]) String() string {
	return fmt.Sprintf("C%0.2x", DataEffect(e))
}

func (e RowJump[TPeriod]) RowEnd(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.SetRow(index.Row(e), true)
}

func (e RowJump[TPeriod]) TraceData() string {
	return e.String()
}
