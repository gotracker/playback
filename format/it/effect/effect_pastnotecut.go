package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/note"
)

// PastNoteCut defines a past note cut effect
type PastNoteCut channel.DataEffect // 'S70'

// Start triggers on the first tick, but before the Tick() function is called
func (e PastNoteCut) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.DoPastNoteEffect(note.ActionCut)
	return nil
}

func (e PastNoteCut) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
