package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetTempo defines a set tempo effect
type SetTempo ChannelCommand // 'T'

func (e SetTempo) String() string {
	return fmt.Sprintf("T%0.2x", DataEffect(e))
}

func (e SetTempo) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	switch DataEffect(e >> 4) {
	case 0: // decrease tempo
		if tick != 0 {
			mem, err := machine.GetChannelMemory[*Memory](m, ch)
			if err != nil {
				return err
			}

			val := int(mem.TempoDecrease(DataEffect(e & 0x0F)))
			if err := m.SlideBPM(-val); err != nil {
				return err
			}
		}
	case 1: // increase tempo
		if tick != 0 {
			mem, err := machine.GetChannelMemory[*Memory](m, ch)
			if err != nil {
				return err
			}

			val := int(mem.TempoIncrease(DataEffect(e & 0x0F)))
			if err := m.SlideBPM(val); err != nil {
				return err
			}
		}
	default:
		if tick == 0 {
			if err := m.SetBPM(int(e)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e SetTempo) TraceData() string {
	return e.String()
}
