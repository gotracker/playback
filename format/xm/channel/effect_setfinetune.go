package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetFinetune defines a mod-style set finetune effect
type SetFinetune[TPeriod period.Period] DataEffect // 'E5x'

func (e SetFinetune[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e SetFinetune[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	inst, err := m.GetChannelInstrument(ch)
	if err != nil {
		return err
	}

	if inst != nil {
		ft := (note.Finetune(e&0x0F) - 8) * 4
		inst.SetFinetune(ft)
	}
	return nil
}

func (e SetFinetune[TPeriod]) TraceData() string {
	return e.String()
}
