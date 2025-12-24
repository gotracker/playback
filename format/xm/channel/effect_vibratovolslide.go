package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
)

// VibratoVolumeSlide defines a combination vibrato and volume slide effect
type VibratoVolumeSlide[TPeriod period.Period] struct { // '6'
	playback.CombinedEffect[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]
}

// NewVibratoVolumeSlide creates a new VibratoVolumeSlide object
func NewVibratoVolumeSlide[TPeriod period.Period](val DataEffect) VibratoVolumeSlide[TPeriod] {
	vvs := VibratoVolumeSlide[TPeriod]{}
	vvs.Effects = append(vvs.Effects, VolumeSlide[TPeriod](val), Vibrato[TPeriod](0x00))
	return vvs
}

func (e VibratoVolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("6%0.2x", DataEffect(any(e.Effects[0]).(VolumeSlide[TPeriod])))
}

func (e VibratoVolumeSlide[TPeriod]) TraceData() string {
	return e.String()
}
