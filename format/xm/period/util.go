package period

import (
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

const (
	// DefaultC2Spd is the default C2SPD for XM samples
	DefaultC2Spd = 8363
	c2Period     = 1712

	floatDefaultC2Spd = float32(DefaultC2Spd)

	// XMBaseClock is the base clock speed of xm files
	XMBaseClock period.Frequency = DefaultC2Spd * c2Period

	NotesPerOctave        = 12
	SlideFinesPerSemitone = 4
	SemitonesPerNote      = 16
	SlideFinesPerNote     = SlideFinesPerSemitone * SemitonesPerNote
	SlideFinesPerOctave   = SlideFinesPerNote * NotesPerOctave

	C4SlideFines = 4 * SlideFinesPerOctave
)

var semitonePeriodTable = [...]float32{27392, 25856, 24384, 23040, 21696, 20480, 19328, 18240, 17216, 16256, 15360, 14496}

// CalcSemitonePeriod calculates the semitone period for it notes
func CalcSemitonePeriod[TPeriod period.Period](semi note.Semitone, ft note.Finetune, c2spd period.Frequency) *TPeriod {
	if semi == note.UnchangedSemitone {
		panic("how?")
	}

	var result TPeriod
	switch p := any(&result).(type) {
	case *Linear:
		nft := int(semi)*SlideFinesPerNote + int(ft)
		p.Finetune = note.Finetune(nft)
		p.CommonRate = c2spd
	case *Amiga:
		key := int(semi.Key())
		octave := uint32(semi.Octave())

		if key >= len(semitonePeriodTable) {
			return nil
		}

		if c2spd == 0 {
			c2spd = period.Frequency(DefaultC2Spd)
		}

		if ft != 0 {
			c2spd = CalcFinetuneC2Spd[Amiga](c2spd, ft)
		}

		p.Amiga = period.Amiga(float64(floatDefaultC2Spd*semitonePeriodTable[key]) / float64(uint32(c2spd)<<octave))
	default:
		return nil
	}

	return &result
}

// CalcFinetuneC2Spd calculates a new C2SPD after a finetune adjustment
func CalcFinetuneC2Spd[TPeriod period.Period](c2spd period.Frequency, finetune note.Finetune) period.Frequency {
	if finetune == 0 {
		return c2spd
	}

	nft := C4SlideFines + int(finetune)
	p := CalcSemitonePeriod[TPeriod](note.Semitone(nft/SlideFinesPerNote), note.Finetune(nft%SlideFinesPerNote), c2spd)
	return any(p).(period.Period).GetFrequency()
}

// FrequencyFromSemitone returns the frequency from the semitone (and c2spd)
func FrequencyFromSemitone[TPeriod period.Period](semitone note.Semitone, c2spd period.Frequency) float32 {
	p := CalcSemitonePeriod[TPeriod](semitone, 0, c2spd)
	return float32(any(p).(period.Period).GetFrequency())
}
