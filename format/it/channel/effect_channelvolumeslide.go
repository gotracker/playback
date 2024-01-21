package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// ChannelVolumeSlide defines a set channel volume effect
type ChannelVolumeSlide[TPeriod period.Period] DataEffect // 'Nxy'

func (e ChannelVolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("N%0.2x", DataEffect(e))
}

func (e ChannelVolumeSlide[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	x, y := mem.ChannelVolumeSlide(DataEffect(e))
	switch {
	case y == 0x0 && x != 0xF:
		// slide up
		return m.SlideChannelMixingVolume(ch, 1, float32(x))
	case y != 0xF && x == 0x0:
		// slide down
		return m.SlideChannelMixingVolume(ch, 1, -float32(y))
	default:
		return nil
	}
}

func (e ChannelVolumeSlide[TPeriod]) TraceData() string {
	return e.String()
}
