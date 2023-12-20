package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

// NewNoteActionNoteContinue defines a NewNoteAction: Note Continue effect
type NewNoteActionNoteContinue[TPeriod period.Period] channel.DataEffect // 'S74'

// Start triggers on the first tick, but before the Tick() function is called
func (e NewNoteActionNoteContinue[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.SetNewNoteAction(note.ActionContinue)
	return nil
}

func (e NewNoteActionNoteContinue[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
