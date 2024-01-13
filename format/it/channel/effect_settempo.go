package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetTempo defines a set tempo effect
type SetTempo[TPeriod period.Period] DataEffect // 'T'

func (e SetTempo[TPeriod]) String() string {
	return fmt.Sprintf("T%0.2x", DataEffect(e))
}

func (e SetTempo[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int) error {
	switch DataEffect(e >> 4) {
	case 0: // decrease BPM
		if tick != 0 {
			mem, err := machine.GetChannelMemory[*Memory](m, ch)
			if err != nil {
				return err
			}
			val := int(mem.TempoDecrease(DataEffect(e & 0x0F)))
			return m.SlideBPM(-val)
		}
	case 1: // increase BPM
		if tick != 0 {
			mem, err := machine.GetChannelMemory[*Memory](m, ch)
			if err != nil {
				return err
			}
			val := int(mem.TempoIncrease(DataEffect(e & 0x0F)))
			return m.SlideBPM(val)
		}
	default:
		if tick == 0 {
			return m.SetBPM(int(e))
		}
	}
	return nil
}

func (e SetTempo[TPeriod]) TraceData() string {
	return e.String()
}
