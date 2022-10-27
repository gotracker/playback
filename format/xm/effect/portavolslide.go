package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// PortaVolumeSlide defines a portamento-to-note combined with a volume slide effect
type PortaVolumeSlide struct { // '5'
	playback.CombinedEffect[channel.State]
}

// NewPortaVolumeSlide creates a new PortaVolumeSlide object
func NewPortaVolumeSlide(val channel.DataEffect) PortaVolumeSlide {
	pvs := PortaVolumeSlide{}
	pvs.Effects = append(pvs.Effects, VolumeSlide(val), PortaToNote(0x00))
	return pvs
}

func (e PortaVolumeSlide) String() string {
	return fmt.Sprintf("5%0.2x", channel.DataEffect(e.Effects[0].(VolumeSlide)))
}
