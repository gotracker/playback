package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PastNoteOff defines a past note off effect
type PastNoteOff[TPeriod period.Period] DataEffect // 'S71'

func (e PastNoteOff[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e PastNoteOff[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.DoChannelPastNoteEffect(ch, note.ActionRelease)
}

func (e PastNoteOff[TPeriod]) TraceData() string {
	return e.String()
}
