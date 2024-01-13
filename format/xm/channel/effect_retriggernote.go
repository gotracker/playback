package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// RetriggerNote defines a retriggering effect
type RetriggerNote[TPeriod period.Period] DataEffect // 'E9x'

func (e RetriggerNote[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e RetriggerNote[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	y := DataEffect(e) & 0x0F
	return m.SetChannelNoteAction(ch, note.ActionRetrigger, int(y))
}

func (e RetriggerNote[TPeriod]) TraceData() string {
	return e.String()
}
