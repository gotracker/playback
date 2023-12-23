package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

// NewNoteActionNoteCut defines a NewNoteAction: Note Cut effect
type NewNoteActionNoteCut[TPeriod period.Period] DataEffect // 'S73'

// Start triggers on the first tick, but before the Tick() function is called
func (e NewNoteActionNoteCut[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.SetNewNoteAction(note.ActionCut)
	return nil
}

func (e NewNoteActionNoteCut[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
