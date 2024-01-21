package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetSpeed defines a set speed effect
type SetSpeed[TPeriod period.Period] DataEffect // 'F'

func (e SetSpeed[TPeriod]) String() string {
	return fmt.Sprintf("F%0.2x", DataEffect(e))
}

func (e SetSpeed[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	if e == 0 {
		return nil
	}
	return m.SetTempo(int(e))
}

func (e SetSpeed[TPeriod]) TraceData() string {
	return e.String()
}
