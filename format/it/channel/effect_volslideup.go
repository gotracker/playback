package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// VolumeSlideUp defines a volume slide up effect
type VolumeSlideUp[TPeriod period.Period] DataEffect // 'D'

func (e VolumeSlideUp[TPeriod]) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}

func (e VolumeSlideUp[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x, _ := mem.VolumeSlide(DataEffect(e))
	return m.SlideChannelVolume(ch, 1, float32(x))
}

func (e VolumeSlideUp[TPeriod]) TraceData() string {
	return e.String()
}

//====================================================

// VolChanVolumeSlideUp defines a volume slide up effect (from the volume channel)
type VolChanVolumeSlideUp[TPeriod period.Period] DataEffect // 'd'

func (e VolChanVolumeSlideUp[TPeriod]) String() string {
	return fmt.Sprintf("d%x0", DataEffect(e))
}

func (e VolChanVolumeSlideUp[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x := mem.VolChanVolumeSlide(DataEffect(e))
	return m.SlideChannelVolume(ch, 1, float32(x))
}

func (e VolChanVolumeSlideUp[TPeriod]) TraceData() string {
	return e.String()
}
