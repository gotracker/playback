package effect

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/oscillator"

	"github.com/gotracker/playback/format/xm/channel"
	effectIntf "github.com/gotracker/playback/format/xm/effect/intf"
	xmPeriod "github.com/gotracker/playback/format/xm/period"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/note"
	"github.com/heucuva/comparison"
)

func doVolSlide[TPeriod period.Period](cs playback.Channel[TPeriod, channel.Memory], delta float32, multiplier float32) error {
	av := cs.GetActiveVolume()
	v := xmVolume.ToVolumeXM(av)
	vol := int16((float32(v) + delta) * multiplier)
	if vol >= 0x40 {
		vol = 0x40
	}
	if vol < 0x00 {
		vol = 0x00
	}
	v = xmVolume.XmVolume(channel.DataEffect(vol))
	cs.SetActiveVolume(v.Volume())
	return nil
}

func doGlobalVolSlide(m effectIntf.XM, delta float32, multiplier float32) error {
	gv := m.GetGlobalVolume()
	v := xmVolume.ToVolumeXM(gv)
	vol := int16((float32(v) + delta) * multiplier)
	if vol >= 0x40 {
		vol = 0x40
	}
	if vol < 0x00 {
		vol = 0x00
	}
	v = xmVolume.XmVolume(channel.DataEffect(vol))
	m.SetGlobalVolume(v.Volume())
	return nil
}

func doPortaByDeltaAmiga(cs playback.Channel[xmPeriod.Amiga, channel.Memory], delta int) error {
	cur := cs.GetPeriod()
	if cur == nil {
		return nil
	}

	d := period.PeriodDelta(delta)
	cur = cur.Add(d)
	cs.SetPeriod(cur)
	return nil
}

func doPortaByDeltaLinear(cs playback.Channel[xmPeriod.Linear, channel.Memory], delta int) error {
	cur := cs.GetPeriod()
	if cur == nil {
		return nil
	}

	finetune := period.PeriodDelta(delta)
	cur = cur.Add(finetune)
	cs.SetPeriod(cur)
	return nil
}

func doPortaUp[TPeriod period.Period](cs playback.Channel[TPeriod, channel.Memory], amount float32, multiplier float32) error {
	delta := int(amount * multiplier)
	switch csp := any(cs).(type) {
	case playback.Channel[xmPeriod.Linear, channel.Memory]:
		return doPortaByDeltaLinear(csp, delta)
	case playback.Channel[xmPeriod.Amiga, channel.Memory]:
		return doPortaByDeltaAmiga(csp, -delta)
	default:
		panic("unhandled channel type")
	}
}

func doPortaUpToNote[TPeriod period.Period](cs playback.Channel[TPeriod, channel.Memory], amount float32, multiplier float32, target *TPeriod) error {
	if err := doPortaUp(cs, amount, multiplier); err != nil {
		return err
	}
	if cur := cs.GetPeriod(); period.ComparePeriods(cur, target) == comparison.SpaceshipLeftGreater {
		cs.SetPeriod(target)
	}
	return nil
}

func doPortaDown[TPeriod period.Period](cs playback.Channel[TPeriod, channel.Memory], amount float32, multiplier float32) error {
	delta := int(amount * multiplier)
	switch csp := any(cs).(type) {
	case playback.Channel[xmPeriod.Linear, channel.Memory]:
		return doPortaByDeltaLinear(csp, -delta)
	case playback.Channel[xmPeriod.Amiga, channel.Memory]:
		return doPortaByDeltaAmiga(csp, delta)
	default:
		panic("unhandled channel type")
	}
}

func doPortaDownToNote[TPeriod period.Period](cs playback.Channel[TPeriod, channel.Memory], amount float32, multiplier float32, target *TPeriod) error {
	if err := doPortaDown(cs, amount, multiplier); err != nil {
		return err
	}
	if cur := cs.GetPeriod(); period.ComparePeriods(cur, target) == comparison.SpaceshipRightGreater {
		cs.SetPeriod(target)
	}
	return nil
}

func doVibrato[TPeriod period.Period](cs playback.Channel[TPeriod, channel.Memory], currentTick int, speed channel.DataEffect, depth channel.DataEffect, multiplier float32) error {
	mem := cs.GetMemory()
	vib := calculateWaveTable(cs, currentTick, speed, depth, multiplier, mem.VibratoOscillator())
	delta := period.PeriodDelta(vib)
	cs.SetPeriodDelta(delta)
	return nil
}

func doTremor[TPeriod period.Period](cs playback.Channel[TPeriod, channel.Memory], currentTick int, onTicks int, offTicks int) error {
	mem := cs.GetMemory()
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
	cs.SetVolumeActive(tremor.IsActive())
	return nil
}

func doArpeggio[TPeriod period.Period](cs playback.Channel[TPeriod, channel.Memory], currentTick int, arpSemitoneADelta int8, arpSemitoneBDelta int8) error {
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
	volSlideTwoThirdsTable = [...]xmVolume.XmVolume{
		0, 0, 1, 1, 2, 3, 3, 4, 5, 5, 6, 6, 7, 8, 8, 9,
		10, 10, 11, 11, 12, 13, 13, 14, 15, 15, 16, 16, 17, 18, 18, 19,
		20, 20, 21, 21, 22, 23, 23, 24, 25, 25, 26, 26, 27, 28, 28, 29,
		30, 30, 31, 31, 32, 33, 33, 34, 35, 35, 36, 36, 37, 38, 38, 39,
	}
)

func doVolSlideTwoThirds[TPeriod period.Period](cs playback.Channel[TPeriod, channel.Memory]) error {
	vol := xmVolume.ToVolumeXM(cs.GetActiveVolume())
	if vol >= 64 {
		vol = 63
	}

	v := volSlideTwoThirdsTable[vol]
	if v >= 0x40 {
		v = 0x40
	}

	cs.SetActiveVolume(v.Volume())
	return nil
}

func doTremolo[TPeriod period.Period](cs playback.Channel[TPeriod, channel.Memory], currentTick int, speed channel.DataEffect, depth channel.DataEffect, multiplier float32) error {
	mem := cs.GetMemory()
	delta := calculateWaveTable(cs, currentTick, speed, depth, multiplier, mem.TremoloOscillator())
	return doVolSlide(cs, delta, 1.0)
}

func calculateWaveTable[TPeriod period.Period](cs playback.Channel[TPeriod, channel.Memory], currentTick int, speed channel.DataEffect, depth channel.DataEffect, multiplier float32, o oscillator.Oscillator) float32 {
	delta := o.GetWave(float32(depth) * multiplier)
	o.Advance(int(speed))
	return delta
}
