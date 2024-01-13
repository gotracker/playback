package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// OrderJump defines an order jump effect
type OrderJump ChannelCommand // 'B'

func (e OrderJump) String() string {
	return fmt.Sprintf("B%0.2x", DataEffect(e))
}

func (e OrderJump) RowEnd(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	o := index.Order(e)
	return m.SetOrder(o)
}

func (e OrderJump) TraceData() string {
	return e.String()
}
