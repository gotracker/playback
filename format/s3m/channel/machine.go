package channel

import (
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

func withOscillatorDo(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], speed int, depth float32, osc machine.Oscillator, fn func(value float32) error) error {
	value, err := m.GetNextChannelWavetableValue(ch, speed, depth, machine.OscillatorVibrato)
	if err != nil {
		return err
	}

	return fn(value)
}

func doPortaUp(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], amount float32, multiplier float32) error {
	cur, err := m.GetChannelPeriod(ch)
	if err != nil {
		return err
	}
	if cur.IsInvalid() {
		return nil
	}

	delta := int(amount * multiplier)
	cur = cur.PortaUp(delta)
	return m.SetChannelPeriod(ch, cur)
}

func doPortaDown(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], amount float32, multiplier float32) error {
	cur, err := m.GetChannelPeriod(ch)
	if err != nil {
		return err
	}
	if cur.IsInvalid() {
		return nil
	}

	delta := int(amount * multiplier)
	cur = cur.PortaDown(delta)
	return m.SetChannelPeriod(ch, cur)
}

func doArpeggio(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int, arpSemitoneADelta, arpSemitoneBDelta int8) error {
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

func doTremor(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], onTicks int, offTicks int) error {
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
