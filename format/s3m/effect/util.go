package effect

import (
	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/oscillator"

	"github.com/gotracker/playback/format/s3m/channel"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/note"
	"github.com/heucuva/comparison"
)

func doVolSlide(cs *channel.State, delta float32, multiplier float32) error {
	av := cs.GetActiveVolume()
	v := s3mVolume.VolumeToS3M(av)
	vol := int16((float32(v) + delta) * multiplier)
	if vol >= 64 {
		vol = 63
	}
	if vol < 0 {
		vol = 0
	}
	sv := s3mfile.Volume(channel.DataEffect(vol))
	nv := s3mVolume.VolumeFromS3M(sv)
	cs.SetActiveVolume(nv)
	return nil
}

func doPortaUp(cs *channel.State, amount float32, multiplier float32) error {
	delta := int(amount * multiplier)
	cs.DoPortaByDelta(-delta)
	return nil
}

func doPortaUpToNote(cs *channel.State, amount float32, multiplier float32, target period.Period) error {
	if err := doPortaUp(cs, amount, multiplier); err != nil {
		return err
	}
	if cur := cs.GetPeriod(); period.ComparePeriods(cur, target) == comparison.SpaceshipLeftGreater {
		cs.SetPeriod(target)
	}
	return nil
}

func doPortaDown(cs *channel.State, amount float32, multiplier float32) error {
	delta := int(amount * multiplier)
	cs.DoPortaByDelta(delta)
	return nil
}

func doPortaDownToNote(cs *channel.State, amount float32, multiplier float32, target period.Period) error {
	if err := doPortaDown(cs, amount, multiplier); err != nil {
		return err
	}
	if cur := cs.GetPeriod(); period.ComparePeriods(cur, target) == comparison.SpaceshipRightGreater {
		cs.SetPeriod(target)
	}
	return nil
}

func doVibrato(cs *channel.State, currentTick int, speed channel.DataEffect, depth channel.DataEffect, multiplier float32) error {
	mem := cs.GetMemory()
	delta := calculateWaveTable(cs, currentTick, channel.DataEffect(speed), channel.DataEffect(depth), multiplier, mem.VibratoOscillator())
	cs.SetPeriodDelta(period.PeriodDelta(delta))
	return nil
}

func doTremor(cs *channel.State, currentTick int, onTicks int, offTicks int) error {
	mem := cs.GetMemory()
	tremor := mem.TremorMem()
	if tremor.IsActive() {
		if tremor.Advance() > onTicks {
			tremor.ToggleAndReset()
		}
	} else {
		if tremor.Advance() > offTicks {
			tremor.ToggleAndReset()
		}
	}
	cs.SetVolumeActive(tremor.IsActive())
	return nil
}

func doArpeggio(cs *channel.State, currentTick int, arpSemitoneADelta int8, arpSemitoneBDelta int8) error {
	ns := cs.GetNoteSemitone()
	var arpSemitoneTarget note.Semitone
	switch currentTick % 3 {
	case 0:
		arpSemitoneTarget = ns
	case 1:
		arpSemitoneTarget = note.Semitone(int8(ns) + arpSemitoneADelta)
	case 2:
		arpSemitoneTarget = note.Semitone(int8(ns) + arpSemitoneBDelta)
	}
	cs.SetOverrideSemitone(arpSemitoneTarget)
	cs.SetTargetPos(cs.GetPos())
	return nil
}

var (
	volSlideTwoThirdsTable = [...]s3mfile.Volume{
		0, 0, 1, 1, 2, 3, 3, 4, 5, 5, 6, 6, 7, 8, 8, 9,
		10, 10, 11, 11, 12, 13, 13, 14, 15, 15, 16, 16, 17, 18, 18, 19,
		20, 20, 21, 21, 22, 23, 23, 24, 25, 25, 26, 26, 27, 28, 28, 29,
		30, 30, 31, 31, 32, 33, 33, 34, 35, 35, 36, 36, 37, 38, 38, 39,
	}
)

func doVolSlideTwoThirds(cs *channel.State) error {
	vol := s3mVolume.VolumeToS3M(cs.GetActiveVolume())
	if vol >= 64 {
		vol = 63
	}
	cs.SetActiveVolume(s3mVolume.VolumeFromS3M(volSlideTwoThirdsTable[vol]))
	return nil
}

func doTremolo(cs *channel.State, currentTick int, speed channel.DataEffect, depth channel.DataEffect, multiplier float32) error {
	mem := cs.GetMemory()
	delta := calculateWaveTable(cs, currentTick, speed, depth, multiplier, mem.TremoloOscillator())
	return doVolSlide(cs, delta, 1.0)
}

func calculateWaveTable(cs *channel.State, currentTick int, speed channel.DataEffect, depth channel.DataEffect, multiplier float32, o oscillator.Oscillator) float32 {
	delta := o.GetWave(float32(depth)) * multiplier
	o.Advance(int(speed))
	return delta
}
