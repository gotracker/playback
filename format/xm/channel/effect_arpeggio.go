package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// Arpeggio defines an arpeggio effect
type Arpeggio[TPeriod period.Period] DataEffect // '0'

func (e Arpeggio[TPeriod]) String() string {
	return fmt.Sprintf("0%0.2x", DataEffect(e))
}

func (e Arpeggio[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	xy := DataEffect(e)
	if xy == 0 {
		return nil
	}

	x, y := int8(xy>>4), int8(xy&0x0f)
	return doArpeggio[TPeriod](ch, m, tick, x, y)
}

func (e Arpeggio[TPeriod]) TraceData() string {
	return e.String()
}
