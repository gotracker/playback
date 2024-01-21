package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FineVolumeSlideUp defines a volume slide effect
type FineVolumeSlideUp[TPeriod period.Period] DataEffect // 'EAx'

func (e FineVolumeSlideUp[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e FineVolumeSlideUp[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}
	y := mem.FineVolumeSlideUp(DataEffect(e)) & 0x0F

	if tick != 0 {
		return nil
	}

	return m.SlideChannelVolume(ch, 1.0, float32(y))
}

func (e FineVolumeSlideUp[TPeriod]) TraceData() string {
	return e.String()
}
