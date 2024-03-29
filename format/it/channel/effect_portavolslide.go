package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/period"
)

// PortaVolumeSlide defines a portamento-to-note combined with a volume slide effect
type PortaVolumeSlide[TPeriod period.Period] struct { // 'L'
	playback.CombinedEffect[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning, *Memory, Data[TPeriod]]
}

// NewPortaVolumeSlide creates a new PortaVolumeSlide object
func NewPortaVolumeSlide[TPeriod period.Period](mem *Memory, cd Command, val DataEffect) PortaVolumeSlide[TPeriod] {
	pvs := PortaVolumeSlide[TPeriod]{}
	vs := volumeSlideFactory[TPeriod](mem, cd, val)
	pvs.Effects = append(pvs.Effects, vs, PortaToNote[TPeriod](0x00))
	return pvs
}

func (e PortaVolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("L%0.2x", any(e.Effects[0]).(DataEffect))
}

func (e PortaVolumeSlide[TPeriod]) TraceData() string {
	return e.String()
}
