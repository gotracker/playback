package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// PortaVolumeSlide defines a portamento-to-note combined with a volume slide effect
type PortaVolumeSlide[TPeriod period.Period] struct { // 'L'
	playback.CombinedEffect[TPeriod, channel.Memory]
}

// NewPortaVolumeSlide creates a new PortaVolumeSlide object
func NewPortaVolumeSlide[TPeriod period.Period](mem *channel.Memory, cd channel.Command, val channel.DataEffect) PortaVolumeSlide[TPeriod] {
	pvs := PortaVolumeSlide[TPeriod]{}
	vs := volumeSlideFactory[TPeriod](mem, cd, val)
	pvs.Effects = append(pvs.Effects, vs, PortaToNote[TPeriod](0x00))
	return pvs
}

func (e PortaVolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("L%0.2x", e.Effects[0].(channel.DataEffect))
}
