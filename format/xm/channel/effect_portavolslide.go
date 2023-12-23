package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// PortaVolumeSlide defines a portamento-to-note combined with a volume slide effect
type PortaVolumeSlide[TPeriod period.Period] struct { // '5'
	playback.CombinedEffect[TPeriod, Memory, Data]
}

// NewPortaVolumeSlide creates a new PortaVolumeSlide object
func NewPortaVolumeSlide[TPeriod period.Period](val DataEffect) PortaVolumeSlide[TPeriod] {
	pvs := PortaVolumeSlide[TPeriod]{}
	pvs.Effects = append(pvs.Effects, VolumeSlide[TPeriod](val), PortaToNote[TPeriod](0x00))
	return pvs
}

func (e PortaVolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("5%0.2x", DataEffect(e.Effects[0].(VolumeSlide[TPeriod])))
}
