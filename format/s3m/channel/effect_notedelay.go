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

// NoteDelay defines a note delay effect
type NoteDelay ChannelCommand // 'SDx'

func (e NoteDelay) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e NoteDelay) RowStart(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	tick := int(DataEffect(e) & 0x0F)
	return m.SetChannelNoteAction(ch, note.ActionRetrigger, tick)
}

func (e NoteDelay) TraceData() string {
	return e.String()
}
