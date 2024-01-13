package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetSpeed defines a set speed effect
type SetSpeed[TPeriod period.Period] DataEffect // 'A'

func (e SetSpeed[TPeriod]) String() string {
	return fmt.Sprintf("A%0.2x", DataEffect(e))
}

func (e SetSpeed[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	return m.SetTempo(int(e))
}

func (e SetSpeed[TPeriod]) TraceData() string {
	return e.String()
}
