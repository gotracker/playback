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

// NewNoteActionNoteOff defines a NewNoteAction: Note Off effect
type NewNoteActionNoteOff[TPeriod period.Period] DataEffect // 'S75'

func (e NewNoteActionNoteOff[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e NewNoteActionNoteOff[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.SetChannelNewNoteAction(ch, note.ActionRelease)
}

func (e NewNoteActionNoteOff[TPeriod]) TraceData() string {
	return e.String()
}
