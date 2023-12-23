package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

// PastNoteFade defines a past note fadeout effect
type PastNoteFade[TPeriod period.Period] DataEffect // 'S72'

// Start triggers on the first tick, but before the Tick() function is called
func (e PastNoteFade[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.DoPastNoteEffect(note.ActionFadeout)
	return nil
}

func (e PastNoteFade[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
