package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetCoarsePanPosition defines a set pan position effect
type SetCoarsePanPosition[TPeriod period.Period] DataEffect // 'E8x'

func (e SetCoarsePanPosition[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e SetCoarsePanPosition[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	pan := xmPanning.Panning((e & 0x0F) << 4)
	return m.SetChannelPan(ch, pan)
}

func (e SetCoarsePanPosition[TPeriod]) TraceData() string {
	return e.String()
}
