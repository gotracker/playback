package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// Tremor defines a tremor effect
type Tremor[TPeriod period.Period] DataEffect // 'T'

func (e Tremor[TPeriod]) String() string {
	return fmt.Sprintf("T%0.2x", DataEffect(e))
}

func (e Tremor[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x, y := mem.Tremor(DataEffect(e))
	return doTremor(ch, m, int(x)+1, int(y)+1)
}

func (e Tremor[TPeriod]) TraceData() string {
	return e.String()
}
