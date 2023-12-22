package period

import (
	"math"

	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

const (
	// DefaultC2Spd is the default C2SPD for IT samples
	DefaultC2Spd = 8363
	// C5Period is the sampler (Amiga-style) period of the C-5 note
	C5Period = 428

	floatDefaultC2Spd = float32(DefaultC2Spd)

	// ITBaseClock is the base clock speed of IT files
	ITBaseClock period.Frequency = DefaultC2Spd * C5Period

	NotesPerOctave        = 12
	SlideFinesPerSemitone = 4
	SemitonesPerNote      = 16
	SlideFinesPerNote     = SlideFinesPerSemitone * SemitonesPerNote
	SlideFinesPerOctave   = SlideFinesPerNote * NotesPerOctave
	C5SlideFines          = 5 * SlideFinesPerOctave
)

var semitonePeriodTable = [...]float32{27392, 25856, 24384, 23040, 21696, 20480, 19328, 18240, 17216, 16256, 15360, 14496}

// CalcSemitonePeriod calculates the semitone period for it notes
func CalcSemitonePeriod[TPeriod period.Period](semi note.Semitone, ft note.Finetune, c2spd period.Frequency) TPeriod {
	if semi == note.UnchangedSemitone {
		panic("how?")
	}
	var result TPeriod
	switch p := any(&result).(type) {
	case *period.Linear:
		nft := int(semi)*SlideFinesPerNote + int(ft)
		p.Finetune = note.Finetune(nft)
		p.CommonRate = c2spd
	case *period.Amiga:
		key := int(semi.Key())
		octave := uint32(semi.Octave())

		if key >= len(semitonePeriodTable) {
			return result
		}

		if c2spd == 0 {
			c2spd = period.Frequency(DefaultC2Spd)
		}

		if ft != 0 {
			c2spd = CalcFinetuneC2Spd[period.Amiga](AmigaConverter, c2spd, ft)
		}

		*p = period.Amiga(float64(floatDefaultC2Spd*semitonePeriodTable[key]) / float64(uint32(c2spd)<<octave))
	default:
	}

	return result
}

// CalcFinetuneC2Spd calculates a new C2SPD after a finetune adjustment
func CalcFinetuneC2Spd[TPeriod period.Period](converter period.PeriodConverter[TPeriod], c2spd period.Frequency, finetune note.Finetune) period.Frequency {
	if finetune == 0 {
		return c2spd
	}

	nft := C5SlideFines + int(finetune)
	p := CalcSemitonePeriod[TPeriod](note.Semitone(nft/SlideFinesPerNote), note.Finetune(nft%SlideFinesPerNote), c2spd)
	return converter.GetFrequency(p)
}

// FrequencyFromSemitone returns the frequency from the semitone (and c2spd)
func FrequencyFromSemitone[TPeriod period.Period](converter period.PeriodConverter[TPeriod], semitone note.Semitone, c2spd period.Frequency) float32 {
	p := CalcSemitonePeriod[TPeriod](semitone, 0, c2spd)
	return float32(converter.GetFrequency(p))
}

// ToAmigaPeriod calculates an amiga period for a linear finetune period
func ToAmigaPeriod(finetunes note.Finetune, c2spd period.Frequency) period.Amiga {
	if finetunes < 0 {
		finetunes = 0
	}
	pow := math.Pow(2, float64(finetunes)/SlideFinesPerOctave)
	linFreq := float64(c2spd) * pow / float64(DefaultC2Spd)

	return period.Amiga(float64(semitonePeriodTable[0]) / linFreq)
}
