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

// NoteDelay defines a note delay effect
type NoteDelay[TPeriod period.Period] DataEffect // 'EDx'

func (e NoteDelay[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e NoteDelay[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	x := DataEffect(e) & 0x0F
	return m.SetChannelNoteAction(ch, note.ActionRetrigger, int(x))
}

func (e NoteDelay[TPeriod]) TraceData() string {
	return e.String()
}
