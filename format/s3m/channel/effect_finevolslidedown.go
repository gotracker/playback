package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FineVolumeSlideDown defines a fine volume slide down effect
type FineVolumeSlideDown ChannelCommand // 'DFy'

func (e FineVolumeSlideDown) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}

func (e FineVolumeSlideDown) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	if tick == 0 {
		return nil
	}

	y := DataEffect(e) & 0x0F

	if y != 0x0F {
		return m.SlideChannelVolume(ch, 1, -float32(y))
	}
	return nil
}

func (e FineVolumeSlideDown) TraceData() string {
	return e.String()
}
