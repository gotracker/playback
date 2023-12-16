package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// VibratoVolumeSlide defines a combination vibrato and volume slide effect
type VibratoVolumeSlide struct { // '6'
	playback.CombinedEffect[channel.Memory]
}

// NewVibratoVolumeSlide creates a new VibratoVolumeSlide object
func NewVibratoVolumeSlide(val channel.DataEffect) VibratoVolumeSlide {
	vvs := VibratoVolumeSlide{}
	vvs.Effects = append(vvs.Effects, VolumeSlide(val), Vibrato(0x00))
	return vvs
}

func (e VibratoVolumeSlide) String() string {
	return fmt.Sprintf("6%0.2x", channel.DataEffect(e.Effects[0].(VolumeSlide)))
}
