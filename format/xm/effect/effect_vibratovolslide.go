package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/period"
)

// VibratoVolumeSlide defines a combination vibrato and volume slide effect
type VibratoVolumeSlide[TPeriod period.Period] struct { // '6'
	playback.CombinedEffect[TPeriod, channel.Memory]
}

// NewVibratoVolumeSlide creates a new VibratoVolumeSlide object
func NewVibratoVolumeSlide[TPeriod period.Period](val channel.DataEffect) VibratoVolumeSlide[TPeriod] {
	vvs := VibratoVolumeSlide[TPeriod]{}
	vvs.Effects = append(vvs.Effects, VolumeSlide[TPeriod](val), Vibrato[TPeriod](0x00))
	return vvs
}

func (e VibratoVolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("6%0.2x", channel.DataEffect(e.Effects[0].(VolumeSlide[TPeriod])))
}
