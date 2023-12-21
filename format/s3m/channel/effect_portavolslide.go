package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
)

// PortaVolumeSlide defines a portamento-to-note combined with a volume slide effect
type PortaVolumeSlide struct { // 'L'
	playback.CombinedEffect[s3mPeriod.Amiga, Memory]
}

// NewPortaVolumeSlide creates a new PortaVolumeSlide object
func NewPortaVolumeSlide(mem *Memory, cd uint8, val DataEffect) PortaVolumeSlide {
	pvs := PortaVolumeSlide{}
	vs := volumeSlideFactory(mem, cd, val)
	pvs.Effects = append(pvs.Effects, vs, PortaToNote(0x00))
	return pvs
}

func (e PortaVolumeSlide) String() string {
	return fmt.Sprintf("L%0.2x", e.Effects[0].(DataEffect))
}
