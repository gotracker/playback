package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FineVolumeSlideUp defines a fine volume slide up effect
type FineVolumeSlideUp ChannelCommand // 'DxF'

func (e FineVolumeSlideUp) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}

func (e FineVolumeSlideUp) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	if tick == 0 {
		return nil
	}

	x := DataEffect(e) >> 4

	if x != 0x0F {
		return m.SlideChannelVolume(ch, 1, float32(x))
	}
	return nil
}

func (e FineVolumeSlideUp) TraceData() string {
	return e.String()
}
