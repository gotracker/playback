package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
)

// VibratoVolumeSlide defines a combination vibrato and volume slide effect
type VibratoVolumeSlide struct { // 'K'
	playback.CombinedEffect[s3mPeriod.Amiga, Memory]
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
