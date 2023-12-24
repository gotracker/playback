package channel

import (
	itfile "github.com/gotracker/goaudiofile/music/tracked/it"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/oscillator"

	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/note"
	"github.com/heucuva/comparison"
)

// IT is an interface to IT effect operations
type IT interface {
	SetTicks(int) error                            // Axx
	SetNextOrder(index.Order) error                // Bxx
	SetNextRow(index.Row) error                    // Cxx
	AddRowTicks(int) error                         // S6x
	SetNextRowWithBacktrack(index.Row, bool) error // SBx
	GetCurrentRow() index.Row                      // SBx
	SetPatternDelay(int) error                     // SEx
	SetTempo(int) error                            // Txx
	IncreaseTempo(int) error                       // Txx
	DecreaseTempo(int) error                       // Txx
	SetGlobalVolume(volume.Volume)                 // Vxx, Wxx
	GetGlobalVolume() volume.Volume                // Vxx, Wxx
	IgnoreUnknownEffect() bool                     // Unhandled
}

func doVolSlide[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], delta float32, multiplier float32) error {
	active := cs.GetActiveState()
	av := active.GetVolume()
	v := itVolume.ToItVolume(av)
	vol := int16((float32(v) + delta) * multiplier)
	if vol >= 0x40 {
		vol = 0x40
	}
	if vol < 0x00 {
		vol = 0x00
	}
	v = itfile.Volume(vol)
	nv := itVolume.FromItVolume(v)
	active.SetVolume(nv)
	return nil
}

func doGlobalVolSlide(m IT, delta float32, multiplier float32) error {
	gv := m.GetGlobalVolume()
	v := itVolume.ToItVolume(gv)
	vol := int16((float32(v) + delta) * multiplier)
	if vol >= 0x40 {
		vol = 0x40
	}
	if vol < 0x00 {
		vol = 0x00
	}
	v = itfile.Volume(vol)
	ngv := itVolume.FromItVolume(v)
	m.SetGlobalVolume(ngv)
	return nil
}

func doPortaUp[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], amount float32, multiplier float32) error {
	active := cs.GetActiveState()
	cur := active.Period
	if cur.IsInvalid() {
		return nil
	}

	delta := int(amount * multiplier)
	cur = period.PortaUp(cur, delta)
	active.Period = cur
	return nil
}

func doPortaUpToNote[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], amount float32, multiplier float32, target TPeriod) error {
	if err := doPortaUp[TPeriod](cs, amount, multiplier); err != nil {
		return err
	}
	if active := cs.GetActiveState(); period.ComparePeriods(active.Period, target) == comparison.SpaceshipLeftGreater {
		active.Period = target
	}
	return nil
}

func doPortaDown[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], amount float32, multiplier float32) error {
	active := cs.GetActiveState()
	cur := active.Period
	if cur.IsInvalid() {
		return nil
	}

	delta := int(amount * multiplier)
	cur = period.PortaDown(cur, delta)
	active.Period = cur
	return nil
}

func doPortaDownToNote[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], amount float32, multiplier float32, target TPeriod) error {
	if err := doPortaDown(cs, amount, multiplier); err != nil {
		return err
	}
	if active := cs.GetActiveState(); period.ComparePeriods(active.Period, target) == comparison.SpaceshipRightGreater {
		active.Period = target
	}
	return nil
}

func doVibrato[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data], currentTick int, speed DataEffect, depth DataEffect, multiplier float32) error {
	mem := cs.GetMemory()
	vib := calculateWaveTable(cs, currentTick, speed, depth, multiplier, mem.VibratoOscillator())
	d := period.Delta(vib)
	cs.SetPeriodDelta(d)
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
	cs.GetTargetState().Pos = cs.GetActiveState().Pos
	return nil
}

var (
	volSlideTwoThirdsTable = [...]DataEffect{
		0, 0, 1, 1, 2, 3, 3, 4, 5, 5, 6, 6, 7, 8, 8, 9,
		10, 10, 11, 11, 12, 13, 13, 14, 15, 15, 16, 16, 17, 18, 18, 19,
		20, 20, 21, 21, 22, 23, 23, 24, 25, 25, 26, 26, 27, 28, 28, 29,
		30, 30, 31, 31, 32, 33, 33, 34, 35, 35, 36, 36, 37, 38, 38, 39,
	}
)

func doVolSlideTwoThirds[TPeriod period.Period](cs playback.Channel[TPeriod, Memory, Data]) error {
	active := cs.GetActiveState()
	vol := itVolume.ToItVolume(active.GetVolume())
	if vol >= 0x10 && vol <= 0x50 {
		vol -= 0x10
		if vol >= 64 {
			vol = 63
		}

		v := volSlideTwoThirdsTable[vol]
		if v >= 0x40 {
			v = 0x40
		}

		vv := itfile.Volume(v)
		active.SetVolume(itVolume.FromItVolume(vv))
	}
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
