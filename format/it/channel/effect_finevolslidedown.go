package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FineVolumeSlideDown defines a fine volume slide down effect
type FineVolumeSlideDown[TPeriod period.Period] DataEffect // 'D'

func (e FineVolumeSlideDown[TPeriod]) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}

func (e FineVolumeSlideDown[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	_, y := mem.VolumeSlide(DataEffect(e))
	if y != 0x0F && tick == 0 {
		return m.SlideChannelVolume(ch, 1.0, -float32(y))
	}
	return nil
}

func (e FineVolumeSlideDown[TPeriod]) TraceData() string {
	return e.String()
}

//====================================================

// VolChanFineVolumeSlideDown defines a fine volume slide down effect (from the volume channel)
type VolChanFineVolumeSlideDown[TPeriod period.Period] DataEffect // 'd'

func (e VolChanFineVolumeSlideDown[TPeriod]) String() string {
	return fmt.Sprintf("dF%x", DataEffect(e))
}

func (e VolChanFineVolumeSlideDown[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	_, y := mem.VolumeSlide(DataEffect(e))
	if tick == 0 {
		return m.SlideChannelVolume(ch, 1.0, -float32(y))
	}
	return nil
}

func (e VolChanFineVolumeSlideDown[TPeriod]) TraceData() string {
	return e.String()
}
