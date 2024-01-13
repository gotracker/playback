package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetFinetune defines a mod-style set finetune effect
type SetFinetune[TPeriod period.Period] DataEffect // 'S2x'

func (e SetFinetune[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e SetFinetune[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	x := DataEffect(e) & 0x0F

	inst, err := m.GetChannelInstrument(ch)
	if err != nil {
		return err
	}

	ft := (note.Finetune(x) - 8) * 4
	inst.SetFinetune(ft)
	return nil
}

func (e SetFinetune[TPeriod]) TraceData() string {
	return e.String()
}
