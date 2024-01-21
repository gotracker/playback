package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// GlobalVolumeSlide defines a global volume slide effect
type GlobalVolumeSlide[TPeriod period.Period] DataEffect // 'W'

func (e GlobalVolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("W%0.2x", DataEffect(e))
}

func (e GlobalVolumeSlide[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	x, y := mem.GlobalVolumeSlide(DataEffect(e))

	if tick == 0 {
		return nil
	}

	if x == 0 {
		// global vol slide down
		return m.SlideGlobalVolume(1, -float32(y))
	} else if y == 0 {
		// global vol slide up
		return m.SlideGlobalVolume(1, float32(x))
	}
	return nil
}

func (e GlobalVolumeSlide[TPeriod]) TraceData() string {
	return e.String()
}
