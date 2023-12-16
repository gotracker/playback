package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
	"github.com/gotracker/playback/note"
)

// NoteDelay defines a note delay effect
type NoteDelay ChannelCommand // 'SDx'

// PreStart triggers when the effect enters onto the channel state
func (e NoteDelay) PreStart(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.SetNotePlayTick(true, note.ActionRetrigger, int(channel.DataEffect(e)&0x0F))
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e NoteDelay) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e NoteDelay) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
