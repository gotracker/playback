package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/voice/types"
)

// Tremolo defines a tremolo effect
type Tremolo[TPeriod period.Period] DataEffect // '7'

func (e Tremolo[TPeriod]) String() string {
	return fmt.Sprintf("7%0.2x", DataEffect(e))
}

func (e Tremolo[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x, y := mem.Tremolo(DataEffect(e))
	// NOTE: JBC - XM updates on tick 0, but MOD does not.
	// Just have to eat this incompatibility, I guess...
	return withOscillatorDo[TPeriod](ch, m, int(x), float32(y)*4, machine.OscillatorTremolo, func(value float32) error {
		return m.SetChannelVolumeDelta(ch, types.VolumeDelta(value))
	})
}

func (e Tremolo[TPeriod]) TraceData() string {
	return e.String()
}
