package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/layout/channel"
	"github.com/gotracker/playback/note"
)

// PastNoteOff defines a past note off effect
type PastNoteOff channel.DataEffect // 'S71'

// Start triggers on the first tick, but before the Tick() function is called
func (e PastNoteOff) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.DoPastNoteEffect(note.ActionRelease)
	return nil
}

func (e PastNoteOff) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
