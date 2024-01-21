package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PatternDelay defines a pattern delay effect
type PatternDelay[TPeriod period.Period] DataEffect // 'SEx'

func (e PatternDelay[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e PatternDelay[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	times := int(DataEffect(e) & 0x0F)
	return m.RowRepeat(times)
}

func (e PatternDelay[TPeriod]) TraceData() string {
	return e.String()
}
