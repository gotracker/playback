package channel

import (
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

func withOscillatorDo[TPeriod period.Period](ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], speed int, depth float32, osc machine.Oscillator, fn func(value float32) error) error {
	value, err := m.GetNextChannelWavetableValue(ch, speed, depth, machine.OscillatorVibrato)
	if err != nil {
		return err
	}

	return fn(value)
}

func doTremor[TPeriod period.Period](ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], onTicks int, offTicks int) error {
	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return err
	}

	tremor := mem.TremorMem()
	if tremor.IsActive() {
		if tremor.Advance() >= onTicks {
			tremor.ToggleAndReset()
		}
	} else {
		if tremor.Advance() >= offTicks {
			tremor.ToggleAndReset()
		}
	}

	return m.SetChannelVolumeActive(ch, tremor.IsActive())
}

func doArpeggio[TPeriod period.Period](ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], tick int, arpSemitoneADelta, arpSemitoneBDelta int8) error {
	switch tick % 3 {
	case 0:
		fallthrough
	default:
		return m.DoChannelArpeggio(ch, 0)
	case 1:
		return m.DoChannelArpeggio(ch, arpSemitoneADelta)
	case 2:
		return m.DoChannelArpeggio(ch, arpSemitoneBDelta)
	}
}
