package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// VibratoVolumeSlide defines a combination vibrato and volume slide effect
type VibratoVolumeSlide[TPeriod period.Period] struct { // 'K'
	playback.CombinedEffect[TPeriod, channel.Memory]
}

// NewVibratoVolumeSlide creates a new VibratoVolumeSlide object
func NewVibratoVolumeSlide[TPeriod period.Period](mem *channel.Memory, cd channel.Command, val channel.DataEffect) VibratoVolumeSlide[TPeriod] {
	vvs := VibratoVolumeSlide[TPeriod]{}
	vs := volumeSlideFactory[TPeriod](mem, cd, val)
	vvs.Effects = append(vvs.Effects, vs, Vibrato[TPeriod](0x00))
	return vvs
}

func (e VibratoVolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("K%0.2x", e.Effects[0].(channel.DataEffect))
}
