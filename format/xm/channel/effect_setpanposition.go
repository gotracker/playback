package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetPanPosition defines a set pan position effect
type SetPanPosition[TPeriod period.Period] DataEffect // '8xx'

func (e SetPanPosition[TPeriod]) String() string {
	return fmt.Sprintf("8%0.2x", DataEffect(e))
}

func (e SetPanPosition[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	if tick != 0 {
		return nil
	}

	xx := uint8(e)
	return m.SetChannelPan(ch, xmPanning.Panning(xx))
}

func (e SetPanPosition[TPeriod]) TraceData() string {
	return e.String()
}
