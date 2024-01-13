package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetFinetune defines a mod-style set finetune effect
type SetFinetune ChannelCommand // 'S2x'

func (e SetFinetune) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e SetFinetune) RowStart(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	x := DataEffect(e) & 0xf

	inst, err := m.GetChannelInstrument(ch)
	if err != nil {
		return err
	}

	if inst != nil {
		ft := (note.Finetune(x) - 8) * 4
		inst.SetFinetune(ft)
	}
	return nil
}

func (e SetFinetune) TraceData() string {
	return e.String()
}
