package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// Vibrato defines a vibrato effect
type Vibrato[TPeriod period.Period] DataEffect // '4'

func (e Vibrato[TPeriod]) String() string {
	return fmt.Sprintf("4%0.2x", DataEffect(e))
}

func (e Vibrato[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x, y := mem.Vibrato(DataEffect(e))

	// NOTE: JBC - XM updates on tick 0, but MOD does not.
	// Just have to eat this incompatibility, I guess...
	return withOscillatorDo[TPeriod](ch, m, int(x), float32(y)*4, machine.OscillatorVibrato, func(value float32) error {
		return m.SetChannelPeriodDelta(ch, period.Delta(value))
	})
}

func (e Vibrato[TPeriod]) TraceData() string {
	return e.String()
}
