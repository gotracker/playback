package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

// PastNoteCut defines a past note cut effect
type PastNoteCut[TPeriod period.Period] channel.DataEffect // 'S70'

// Start triggers on the first tick, but before the Tick() function is called
func (e PastNoteCut[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.DoPastNoteEffect(note.ActionCut)
	return nil
}

func (e PastNoteCut[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
