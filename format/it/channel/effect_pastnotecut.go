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

// PastNoteCut defines a past note cut effect
type PastNoteCut[TPeriod period.Period] DataEffect // 'S70'

func (e PastNoteCut[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e PastNoteCut[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.DoChannelPastNoteEffect(ch, note.ActionCut)
}

func (e PastNoteCut[TPeriod]) TraceData() string {
	return e.String()
}
