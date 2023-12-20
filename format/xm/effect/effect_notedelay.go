package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

// NoteDelay defines a note delay effect
type NoteDelay[TPeriod period.Period] channel.DataEffect // 'EDx'

// PreStart triggers when the effect enters onto the channel state
func (e NoteDelay[TPeriod]) PreStart(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.SetNotePlayTick(true, note.ActionRetrigger, int(channel.DataEffect(e)&0x0F))
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e NoteDelay[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e NoteDelay[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
