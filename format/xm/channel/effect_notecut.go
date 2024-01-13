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

// NoteCut defines a note cut effect
type NoteCut[TPeriod period.Period] DataEffect // 'ECx'

func (e NoteCut[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e NoteCut[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	x := DataEffect(e) & 0x0F
	return m.SetChannelNoteAction(ch, note.ActionCut, int(x))
}

func (e NoteCut[TPeriod]) TraceData() string {
	return e.String()
}
