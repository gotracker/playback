package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// GlobalVolumeSlide defines a global volume slide effect
type GlobalVolumeSlide[TPeriod period.Period] DataEffect // 'H'

func (e GlobalVolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("H%0.2x", DataEffect(e))
}

func (e GlobalVolumeSlide[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x, y := mem.GlobalVolumeSlide(DataEffect(e))

	if tick == 0 {
		return nil
	}

	if x == 0 {
		return m.SlideGlobalVolume(1, -float32(y))
	} else if y == 0 {
		return m.SlideGlobalVolume(1, float32(y))
	}
	return nil
}

func (e GlobalVolumeSlide[TPeriod]) TraceData() string {
	return e.String()
}
