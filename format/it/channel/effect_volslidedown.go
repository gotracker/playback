package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// VolumeSlideDown defines a volume slide down effect
type VolumeSlideDown[TPeriod period.Period] DataEffect // 'D'

func (e VolumeSlideDown[TPeriod]) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}

func (e VolumeSlideDown[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	_, y := mem.VolumeSlide(DataEffect(e))
	return m.SlideChannelVolume(ch, 1.0, -float32(y))
}

func (e VolumeSlideDown[TPeriod]) TraceData() string {
	return e.String()
}

//====================================================

// VolChanVolumeSlideDown defines a volume slide down effect (from the volume channel)
type VolChanVolumeSlideDown[TPeriod period.Period] DataEffect // 'd'

func (e VolChanVolumeSlideDown[TPeriod]) String() string {
	return fmt.Sprintf("d0%x", DataEffect(e))
}

func (e VolChanVolumeSlideDown[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	y := mem.VolChanVolumeSlide(DataEffect(e))
	return m.SlideChannelVolume(ch, 1, -float32(y))
}

func (e VolChanVolumeSlideDown[TPeriod]) TraceData() string {
	return e.String()
}
