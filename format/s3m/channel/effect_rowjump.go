package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// RowJump defines a row jump effect
type RowJump ChannelCommand // 'C'

func (e RowJump) String() string {
	return fmt.Sprintf("C%0.2x", DataEffect(e))
}

func (e RowJump) RowEnd(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	r := DataEffect(e)
	rowIdx := index.Row((r >> 4) * 10)
	rowIdx += index.Row(r & 0xf)

	return m.SetRow(rowIdx, true)
}

func (e RowJump) TraceData() string {
	return e.String()
}
