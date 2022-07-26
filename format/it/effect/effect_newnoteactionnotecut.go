package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/note"
)

// NewNoteActionNoteCut defines a NewNoteAction: Note Cut effect
type NewNoteActionNoteCut channel.DataEffect // 'S73'

// Start triggers on the first tick, but before the Tick() function is called
func (e NewNoteActionNoteCut) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.SetNewNoteAction(note.ActionCut)
	return nil
}

func (e NewNoteActionNoteCut) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
