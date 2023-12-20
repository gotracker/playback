package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

// NewNoteActionNoteOff defines a NewNoteAction: Note Off effect
type NewNoteActionNoteOff[TPeriod period.Period] channel.DataEffect // 'S75'

// Start triggers on the first tick, but before the Tick() function is called
func (e NewNoteActionNoteOff[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.SetNewNoteAction(note.ActionRelease)
	return nil
}

func (e NewNoteActionNoteOff[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
