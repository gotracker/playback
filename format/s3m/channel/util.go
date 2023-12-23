package channel

import (
	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/gomixing/volume"
	"github.com/heucuva/comparison"

	"github.com/gotracker/playback"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/oscillator"
)

type EffectS3M = playback.Effect
type S3MChannel = playback.Channel[period.Amiga, Memory, Data]

// S3M is an interface to S3M effect operations
type S3M interface {
	SetTicks(int) error                            // Axx
	SetNextOrder(index.Order) error                // Bxx
	SetNextRow(index.Row) error                    // Cxx
	SetFilterEnable(bool)                          // S0x
	SetNextRowWithBacktrack(index.Row, bool) error // SBx
	GetCurrentRow() index.Row                      // SBx
	SetPatternDelay(int) error                     // SEx
	AddRowTicks(int) error                         // S6x
	SetTempo(int) error                            // Txx
	IncreaseTempo(int) error                       // Txx
	DecreaseTempo(int) error                       // Txx
	SetGlobalVolume(volume.Volume)                 // Vxx
	IgnoreUnknownEffect() bool                     // Unhandled
}

func doVolSlide(cs S3MChannel, delta float32, multiplier float32) error {
	av := cs.GetActiveVolume()
	v := s3mVolume.VolumeToS3M(av)
	vol := int16((float32(v) + delta) * multiplier)
	if vol >= 64 {
		vol = 63
	}
	if vol < 0 {
		vol = 0
	}
	sv := s3mfile.Volume(vol)
	nv := s3mVolume.VolumeFromS3M(sv)
	cs.SetActiveVolume(nv)
	return nil
}

func doPortaUp(cs S3MChannel, amount float32, multiplier float32) error {
	cur := cs.GetPeriod()
	if cur.IsInvalid() {
		return nil
	}

	delta := int(amount * multiplier)
	cur = cur.PortaUp(delta)
	cs.SetPeriod(cur)
	return nil
}

func doPortaUpToNote(cs S3MChannel, amount float32, multiplier float32, target period.Amiga) error {
	if target.IsInvalid() {
		return nil
	}

	cur := cs.GetPeriod()
	if cur.IsInvalid() {
		return nil
	}

	delta := int(amount * multiplier)
	d := period.Delta(delta)
	cur = cur.Add(d)
	if period.ComparePeriods(cur, target) == comparison.SpaceshipLeftGreater {
		cur = target
	}
	cs.SetPeriod(cur)
	return nil
}

func doPortaDown(cs S3MChannel, amount float32, multiplier float32) error {
	cur := cs.GetPeriod()
	if cur.IsInvalid() {
		return nil
	}

	delta := int(amount * multiplier)
	cur = cur.PortaDown(delta)
	cs.SetPeriod(cur)
	return nil
}

func doPortaDownToNote(cs S3MChannel, amount float32, multiplier float32, target period.Amiga) error {
	cur := cs.GetPeriod()
	if cur.IsInvalid() {
		return nil
	}

	delta := int(amount * multiplier)
	d := period.Delta(-delta)
	cur = cur.Add(d)
	if period.ComparePeriods(cur, target) == comparison.SpaceshipRightGreater {
		cur = target
	}
	cs.SetPeriod(cur)
	return nil
}

func doVibrato(cs S3MChannel, currentTick int, speed DataEffect, depth DataEffect, multiplier float32) error {
	mem := cs.GetMemory()
	delta := calculateWaveTable(cs, currentTick, DataEffect(speed), DataEffect(depth), multiplier, mem.VibratoOscillator())
	d := period.Delta(delta)
	cs.SetPeriodDelta(d)
	return nil
}

func doTremor(cs S3MChannel, currentTick int, onTicks int, offTicks int) error {
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

func doArpeggio(cs S3MChannel, currentTick int, arpSemitoneADelta int8, arpSemitoneBDelta int8) error {
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

func doVolSlideTwoThirds(cs S3MChannel) error {
	vol := s3mVolume.VolumeToS3M(cs.GetActiveVolume())
	if vol >= 64 {
		vol = 63
	}
	cs.SetActiveVolume(s3mVolume.VolumeFromS3M(volSlideTwoThirdsTable[vol]))
	return nil
}

func doTremolo(cs S3MChannel, currentTick int, speed DataEffect, depth DataEffect, multiplier float32) error {
	mem := cs.GetMemory()
	delta := calculateWaveTable(cs, currentTick, speed, depth, multiplier, mem.TremoloOscillator())
	return doVolSlide(cs, delta, 1.0)
}

func calculateWaveTable(cs S3MChannel, currentTick int, speed DataEffect, depth DataEffect, multiplier float32, o oscillator.Oscillator) float32 {
	delta := o.GetWave(float32(depth)) * multiplier
	o.Advance(int(speed))
	return delta
}
