package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetTempo defines a set tempo effect
type SetTempo[TPeriod period.Period] DataEffect // 'F'

func (e SetTempo[TPeriod]) String() string {
	return fmt.Sprintf("F%0.2x", DataEffect(e))
}

func (e SetTempo[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	if e < 0x20 {
		return nil
	}
	return m.SetBPM(int(e))
}

func (e SetTempo[TPeriod]) TraceData() string {
	return e.String()
}
