package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
)

// VibratoVolumeSlide defines a combination vibrato and volume slide effect
type VibratoVolumeSlide struct { // 'K'
	playback.CombinedEffect[channel.State]
}

// NewVibratoVolumeSlide creates a new VibratoVolumeSlide object
func NewVibratoVolumeSlide(mem *channel.Memory, cd uint8, val channel.DataEffect) VibratoVolumeSlide {
	vvs := VibratoVolumeSlide{}
	vs := volumeSlideFactory(mem, cd, val)
	vvs.Effects = append(vvs.Effects, vs, Vibrato(0x00))
	return vvs
}

func (e VibratoVolumeSlide) String() string {
	return fmt.Sprintf("K%0.2x", e.Effects[0].(channel.DataEffect))
}
