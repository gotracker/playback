package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

// PastNoteFade defines a past note fadeout effect
type PastNoteFade[TPeriod period.Period] channel.DataEffect // 'S72'

// Start triggers on the first tick, but before the Tick() function is called
func (e PastNoteFade[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.DoPastNoteEffect(note.ActionFadeout)
	return nil
}

func (e PastNoteFade[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
