package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// RetrigVolumeSlide defines a retriggering volume slide effect
type RetrigVolumeSlide[TPeriod period.Period] DataEffect // 'Q'

func (e RetrigVolumeSlide[TPeriod]) String() string {
	return fmt.Sprintf("Q%0.2x", DataEffect(e))
}

func (e RetrigVolumeSlide[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	x := DataEffect(e) >> 4   // vol slide instruction
	y := DataEffect(e) & 0x0F // number of ticks between retriggers

	if (tick % int(y+1)) != 0 {
		return nil
	}

	if err := m.SetChannelNoteAction(ch, note.ActionRetrigger, tick); err != nil {
		return err
	}

	switch x {
	case 0: // nothing
		fallthrough
	default:

	case 1: // -1
		return m.SlideChannelVolume(ch, 1, -1)
	case 2: // -2
		return m.SlideChannelVolume(ch, 1, -2)
	case 3: // -4
		return m.SlideChannelVolume(ch, 1, -4)
	case 4: // -8
		return m.SlideChannelVolume(ch, 1, -8)
	case 5: // -16
		return m.SlideChannelVolume(ch, 1, -16)
	case 6: // * 2/3
		return m.SlideChannelVolume(ch, 2.0/3.0, 0)
	case 7: // * 1/2
		return m.SlideChannelVolume(ch, 1.0/2.0, 0)
	case 8: // ?
	case 9: // +1
		return m.SlideChannelVolume(ch, 1, 1)
	case 10: // +2
		return m.SlideChannelVolume(ch, 1, 2)
	case 11: // +4
		return m.SlideChannelVolume(ch, 1, 4)
	case 12: // +8
		return m.SlideChannelVolume(ch, 1, 8)
	case 13: // +16
		return m.SlideChannelVolume(ch, 1, 16)
	case 14: // * 3/2
		return m.SlideChannelVolume(ch, 3.0/2.0, 0)
	case 15: // * 2
		return m.SlideChannelVolume(ch, 2, 0)
	}
	return nil
}

func (e RetrigVolumeSlide[TPeriod]) TraceData() string {
	return e.String()
}
