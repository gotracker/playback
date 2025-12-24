package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SampleOffset defines a sample offset effect
type SampleOffset[TPeriod period.Period] DataEffect // '9'

func (e SampleOffset[TPeriod]) String() string {
	return fmt.Sprintf("9%0.2x", DataEffect(e))
}

func (e SampleOffset[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	xx := mem.SampleOffset(DataEffect(e))
	return m.SetChannelPos(ch, sampling.Pos{Pos: int(xx) * 0x100})
}

func (e SampleOffset[TPeriod]) TraceData() string {
	return e.String()
}
