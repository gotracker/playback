package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// OrderJump defines an order jump effect
type OrderJump[TPeriod period.Period] DataEffect // 'B'

func (e OrderJump[TPeriod]) String() string {
	return fmt.Sprintf("B%0.2x", DataEffect(e))
}

func (e OrderJump[TPeriod]) RowEnd(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.SetOrder(index.Order(e))
}

func (e OrderJump[TPeriod]) TraceData() string {
	return e.String()
}
