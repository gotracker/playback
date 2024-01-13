package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// VolumeSlideDown defines a volume slide down effect
type VolumeSlideDown ChannelCommand // 'D0y'

func (e VolumeSlideDown) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}

func (e VolumeSlideDown) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	y := DataEffect(e) & 0x0F
	if mem.Shared.VolSlideEveryFrame || tick != 0 {
		return m.SlideChannelVolume(ch, 1, -float32(y))
	}
	return nil
}

func (e VolumeSlideDown) TraceData() string {
	return e.String()
}
