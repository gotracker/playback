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

// PastNoteFade defines a past note fadeout effect
type PastNoteFade[TPeriod period.Period] DataEffect // 'S72'

func (e PastNoteFade[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e PastNoteFade[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.DoChannelPastNoteEffect(ch, note.ActionFadeout)
}

func (e PastNoteFade[TPeriod]) TraceData() string {
	return e.String()
}
