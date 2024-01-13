package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// VolumeSlideUp defines a volume slide up effect
type VolumeSlideUp ChannelCommand // 'Dx0'

func (e VolumeSlideUp) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}

func (e VolumeSlideUp) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	x := DataEffect(e) >> 4
	if mem.Shared.VolSlideEveryFrame || tick != 0 {
		return m.SlideChannelVolume(ch, 1, float32(x))
	}
	return nil
}

func (e VolumeSlideUp) TraceData() string {
	return e.String()
}
