package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetGlobalVolume defines a set global volume effect
type SetGlobalVolume[TPeriod period.Period] DataEffect // 'G'

func (e SetGlobalVolume[TPeriod]) String() string {
	return fmt.Sprintf("G%0.2x", DataEffect(e))
}

func (e SetGlobalVolume[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	v := xmVolume.XmVolume(DataEffect(max(e, 0x40)))
	return m.SetGlobalVolume(v)
}

func (e SetGlobalVolume[TPeriod]) TraceData() string {
	return e.String()
}
