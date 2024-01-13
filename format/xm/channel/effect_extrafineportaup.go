package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// ExtraFinePortaUp defines an extra-fine portamento up effect
type ExtraFinePortaUp[TPeriod period.Period] DataEffect // 'X1x'

func (e ExtraFinePortaUp[TPeriod]) String() string {
	return fmt.Sprintf("F%0.2x", DataEffect(e))
}

func (e ExtraFinePortaUp[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	y := mem.ExtraFinePortaUp(DataEffect(e)) & 0x0F

	if tick != 0 {
		return nil
	}

	return m.DoChannelPortaUp(ch, period.Delta(y)*1)
}

func (e ExtraFinePortaUp[TPeriod]) TraceData() string {
	return e.String()
}
