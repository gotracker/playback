package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// VibratoVolumeSlide defines a combination vibrato and volume slide effect
type VibratoVolumeSlide struct { // 'K'
	playback.CombinedEffect[period.Amiga, Memory]
}

// NewVibratoVolumeSlide creates a new VibratoVolumeSlide object
func NewVibratoVolumeSlide(mem *Memory, cd uint8, val DataEffect) VibratoVolumeSlide {
	vvs := VibratoVolumeSlide{}
	vs := volumeSlideFactory(mem, cd, val)
	vvs.Effects = append(vvs.Effects, vs, Vibrato(0x00))
	return vvs
}

func (e VibratoVolumeSlide) String() string {
	return fmt.Sprintf("K%0.2x", e.Effects[0].(DataEffect))
}
