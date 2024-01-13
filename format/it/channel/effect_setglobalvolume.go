package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetGlobalVolume defines a set global volume effect
type SetGlobalVolume[TPeriod period.Period] DataEffect // 'V'

func (e SetGlobalVolume[TPeriod]) String() string {
	return fmt.Sprintf("V%0.2x", DataEffect(e))
}

func (e SetGlobalVolume[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	v := max(itVolume.FineVolume(DataEffect(e)), itVolume.MaxItFineVolume)
	return m.SetGlobalVolume(v)
}

func (e SetGlobalVolume[TPeriod]) TraceData() string {
	return e.String()
}
