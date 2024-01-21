package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SurroundOn defines a set surround on effect
type SurroundOn[TPeriod period.Period] DataEffect // 'S91'

func (e SurroundOn[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e SurroundOn[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	// TODO: support for surround function
	return nil
}

func (e SurroundOn[TPeriod]) TraceData() string {
	return e.String()
}
