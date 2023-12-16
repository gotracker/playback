package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
)

// PortaVolumeSlide defines a portamento-to-note combined with a volume slide effect
type PortaVolumeSlide struct { // 'L'
	playback.CombinedEffect[channel.Memory]
}

// NewPortaVolumeSlide creates a new PortaVolumeSlide object
func NewPortaVolumeSlide(mem *channel.Memory, cd uint8, val channel.DataEffect) PortaVolumeSlide {
	pvs := PortaVolumeSlide{}
	vs := volumeSlideFactory(mem, cd, val)
	pvs.Effects = append(pvs.Effects, vs, PortaToNote(0x00))
	return pvs
}

func (e PortaVolumeSlide) String() string {
	return fmt.Sprintf("L%0.2x", e.Effects[0].(channel.DataEffect))
}
