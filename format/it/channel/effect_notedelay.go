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

// NoteDelay defines a note delay effect
type NoteDelay[TPeriod period.Period] DataEffect // 'SDx'

func (e NoteDelay[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e NoteDelay[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.SetChannelNoteAction(ch, note.ActionRetrigger, int(e&0x0F))
}

func (e NoteDelay[TPeriod]) TraceData() string {
	return e.String()
}
