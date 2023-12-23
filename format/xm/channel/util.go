package channel

import (
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/oscillator"

	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/note"
	"github.com/heucuva/comparison"
)

// XM is an interface to XM effect operations
type XM interface {
	SetNextOrder(index.Order) error                // Bxx
	BreakOrder() error                             // Dxx
	SetNextRow(index.Row) error                    // Dxx
	SetNextRowWithBacktrack(index.Row, bool) error // E6x
	GetCurrentRow() index.Row                      // E6x
	SetPatternDelay(int) error                     // EEx
	SetTicks(int) error                            // Fxx
	SetTempo(int) error                            // Fxx
	SetGlobalVolume(volume.Volume)                 // Gxx
	GetGlobalVolume() volume.Volume                // Hxx
	SetEnvelopePosition(int)                       // Lxx
	IgnoreUnknownEffect() bool                     // Unhandled
}

func doVolSlide[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], delta float32, multiplier float32) error {
	av := cs.GetActiveVolume()
	v := xmVolume.ToVolumeXM(av)
	vol := int16((float32(v) + delta) * multiplier)
	if vol >= 0x40 {
		vol = 0x40
	}
	if vol < 0x00 {
		vol = 0x00
	}
	v = xmVolume.XmVolume(DataEffect(vol))
	cs.SetActiveVolume(v.Volume())
	return nil
}

func doGlobalVolSlide(m XM, delta float32, multiplier float32) error {
	gv := m.GetGlobalVolume()
	v := xmVolume.ToVolumeXM(gv)
	vol := int16((float32(v) + delta) * multiplier)
	if vol >= 0x40 {
		vol = 0x40
	}
	if vol < 0x00 {
		vol = 0x00
	}
	v = xmVolume.XmVolume(DataEffect(vol))
	m.SetGlobalVolume(v.Volume())
	return nil
}

func doPortaUp[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], amount float32, multiplier float32) error {
	cur := cs.GetPeriod()
	if cur.IsInvalid() {
		return nil
	}

	delta := int(amount * multiplier)
	cur = period.PortaUp(cur, delta)
	cs.SetPeriod(cur)
	return nil
}

func doPortaUpToNote[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], amount float32, multiplier float32, target TPeriod) error {
	if err := doPortaUp(cs, amount, multiplier); err != nil {
		return err
	}
	if cur := cs.GetPeriod(); period.ComparePeriods(cur, target) == comparison.SpaceshipLeftGreater {
		cs.SetPeriod(target)
	}
	return nil
}

func doPortaDown[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], amount float32, multiplier float32) error {
	cur := cs.GetPeriod()
	if cur.IsInvalid() {
		return nil
	}

	delta := int(amount * multiplier)
	cur = period.PortaDown(cur, delta)
	cs.SetPeriod(cur)
	return nil
}

func doPortaDownToNote[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], amount float32, multiplier float32, target TPeriod) error {
	if err := doPortaDown(cs, amount, multiplier); err != nil {
		return err
	}
	if cur := cs.GetPeriod(); period.ComparePeriods(cur, target) == comparison.SpaceshipRightGreater {
		cs.SetPeriod(target)
	}
	return nil
}

func doVibrato[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], currentTick int, speed DataEffect, depth DataEffect, multiplier float32) error {
	mem := cs.GetMemory()
	vib := calculateWaveTable(cs, currentTick, speed, depth, multiplier, mem.VibratoOscillator())
	delta := period.Delta(vib)
	cs.SetPeriodDelta(delta)
	return nil
}

func doTremor[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], currentTick int, onTicks int, offTicks int) error {
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

func doArpeggio[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], currentTick int, arpSemitoneADelta int8, arpSemitoneBDelta int8) error {
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

func doVolSlideTwoThirds[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data]) error {
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

func doTremolo[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], currentTick int, speed DataEffect, depth DataEffect, multiplier float32) error {
	mem := cs.GetMemory()
	delta := calculateWaveTable(cs, currentTick, speed, depth, multiplier, mem.TremoloOscillator())
	return doVolSlide(cs, delta, 1.0)
}

func calculateWaveTable[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], currentTick int, speed DataEffect, depth DataEffect, multiplier float32, o oscillator.Oscillator) float32 {
	delta := o.GetWave(float32(depth) * multiplier)
	o.Advance(int(speed))
	return delta
}
