package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// OrderJump defines an order jump effect
type OrderJump[TPeriod period.Period] DataEffect // 'B'

func (e OrderJump[TPeriod]) String() string {
	return fmt.Sprintf("B%0.2x", DataEffect(e))
}

func (e OrderJump[TPeriod]) RowEnd(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	return m.SetOrder(index.Order(e))
}

func (e OrderJump[TPeriod]) TraceData() string {
	return e.String()
}
