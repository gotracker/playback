package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// RowJump defines a row jump effect
type RowJump[TPeriod period.Period] DataEffect // 'D'

func (e RowJump[TPeriod]) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}

func (e RowJump[TPeriod]) RowEnd(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	xy := DataEffect(e)
	x, y := xy>>4, xy&0x0f
	row := index.Row(x*10 + y)

	return m.SetRow(row, true)
}

func (e RowJump[TPeriod]) TraceData() string {
	return e.String()
}
