package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetVolume defines a volume slide effect
type SetVolume[TPeriod period.Period] DataEffect // 'C'

func (e SetVolume[TPeriod]) String() string {
	return fmt.Sprintf("C%0.2x", DataEffect(e))
}

func (e SetVolume[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	xx := xmVolume.XmVolume(e)
	return m.SetChannelVolume(ch, xx)
}

func (e SetVolume[TPeriod]) TraceData() string {
	return e.String()
}
