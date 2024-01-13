package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FinePatternDelay defines an fine pattern delay effect
type FinePatternDelay[TPeriod period.Period] DataEffect // 'S6x'

func (e FinePatternDelay[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e FinePatternDelay[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	x := DataEffect(e) & 0x0F
	return m.AddExtraTicks(int(x))
}

func (e FinePatternDelay[TPeriod]) TraceData() string {
	return e.String()
}
