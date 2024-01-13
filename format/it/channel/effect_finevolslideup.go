package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// FineVolumeSlideUp defines a fine volume slide up effect
type FineVolumeSlideUp[TPeriod period.Period] DataEffect // 'D'

func (e FineVolumeSlideUp[TPeriod]) String() string {
	return fmt.Sprintf("D%0.2x", DataEffect(e))
}

func (e FineVolumeSlideUp[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x, _ := mem.VolumeSlide(DataEffect(e))
	if x != 0x0F && tick == 0 {
		return m.SlideChannelVolume(ch, 1.0, float32(x))
	}
	return nil
}

func (e FineVolumeSlideUp[TPeriod]) TraceData() string {
	return e.String()
}

//====================================================

// VolChanFineVolumeSlideUp defines a fine volume slide up effect (from the volume channel)
type VolChanFineVolumeSlideUp[TPeriod period.Period] DataEffect // 'd'

func (e VolChanFineVolumeSlideUp[TPeriod]) String() string {
	return fmt.Sprintf("d%xF", DataEffect(e))
}

func (e VolChanFineVolumeSlideUp[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	x, _ := mem.VolumeSlide(DataEffect(e))
	if tick == 0 {
		return m.SlideChannelVolume(ch, 1.0, float32(x))
	}
	return nil
}

func (e VolChanFineVolumeSlideUp[TPeriod]) TraceData() string {
	return e.String()
}
