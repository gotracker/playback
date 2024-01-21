package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// VolumeSlide defines a volume slide effect
type VolumeSlide[TPeriod period.Period] DataEffect // 'A'

func (e VolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("A%0.2x", DataEffect(e))
}

func (e VolumeSlide[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x, y := mem.VolumeSlide(DataEffect(e))

	if tick == 0 {
		return nil
	}

	if x == 0 {
		// vol slide down
		return m.SlideChannelVolume(ch, 1, -float32(y))
	} else if y == 0 {
		// vol slide up
		return m.SlideChannelVolume(ch, 1, float32(x))
	}

	return nil
}

func (e VolumeSlide[TPeriod]) TraceData() string {
	return e.String()
}
