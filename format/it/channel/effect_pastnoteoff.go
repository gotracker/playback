package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

// PastNoteOff defines a past note off effect
type PastNoteOff[TPeriod period.Period] DataEffect // 'S71'

// Start triggers on the first tick, but before the Tick() function is called
func (e PastNoteOff[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.DoPastNoteEffect(note.ActionRelease)
	return nil
}

func (e PastNoteOff[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
