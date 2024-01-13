package channel

import (
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

func withOscillatorDo[TPeriod period.Period](ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], speed int, depth float32, osc machine.Oscillator, fn func(value float32) error) error {
	value, err := m.GetNextChannelWavetableValue(ch, speed, depth, machine.OscillatorVibrato)
	if err != nil {
		return err
	}

	return fn(value)
}

func doArpeggio[TPeriod period.Period](ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int, arpSemitoneADelta int8, arpSemitoneBDelta int8) error {
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

func doTremor[TPeriod period.Period](ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], onTicks int, offTicks int) error {
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
