package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

// NewNoteActionNoteFade defines a NewNoteAction: Note Fade effect
type NewNoteActionNoteFade[TPeriod period.Period] channel.DataEffect // 'S76'

// Start triggers on the first tick, but before the Tick() function is called
func (e NewNoteActionNoteFade[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.SetNewNoteAction(note.ActionFadeout)
	return nil
}

func (e NewNoteActionNoteFade[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
