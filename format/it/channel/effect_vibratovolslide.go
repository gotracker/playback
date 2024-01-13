package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/period"
)

// VibratoVolumeSlide defines a combination vibrato and volume slide effect
type VibratoVolumeSlide[TPeriod period.Period] struct { // 'K'
	playback.CombinedEffect[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning, *Memory, Data[TPeriod]]
}

// NewVibratoVolumeSlide creates a new VibratoVolumeSlide object
func NewVibratoVolumeSlide[TPeriod period.Period](mem *Memory, cd Command, val DataEffect) VibratoVolumeSlide[TPeriod] {
	vvs := VibratoVolumeSlide[TPeriod]{}
	vs := volumeSlideFactory[TPeriod](mem, cd, val)
	vvs.Effects = append(vvs.Effects, vs, Vibrato[TPeriod](0x00))
	return vvs
}

func (e VibratoVolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("K%0.2x", any(e.Effects[0]).(DataEffect))
}

func (e VibratoVolumeSlide[TPeriod]) TraceData() string {
	return e.String()
}
