package period

import (
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

const (
	// DefaultC4SampleRate is the default c4 sample rate for XM samples
	DefaultC4SampleRate = 8363
	c2Period            = 1712

	floatDefaultC4SampleRate = float32(DefaultC4SampleRate)

	// XMBaseClock is the base clock speed of xm files
	XMBaseClock period.Frequency = DefaultC4SampleRate * c2Period

	NotesPerOctave        = 12
	SlideFinesPerSemitone = 4
	SemitonesPerNote      = 16
	SlideFinesPerNote     = SlideFinesPerSemitone * SemitonesPerNote
	SlideFinesPerOctave   = SlideFinesPerNote * NotesPerOctave

	C4SlideFines = 4 * SlideFinesPerOctave
)

var semitonePeriodTable = [...]float32{27392, 25856, 24384, 23040, 21696, 20480, 19328, 18240, 17216, 16256, 15360, 14496}

// CalcSemitonePeriod calculates the semitone period for it notes
func CalcSemitonePeriod[TPeriod period.Period](semi note.Semitone, ft note.Finetune, c4SampleRate period.Frequency) TPeriod {
	if semi == note.UnchangedSemitone {
		panic("how?")
	}

	var result TPeriod
	switch p := any(&result).(type) {
	case *period.Linear:
		nft := int(semi)*SlideFinesPerNote + int(ft)
		p.Finetune = note.Finetune(nft)
		p.CommonRate = c4SampleRate
	case *period.Amiga:
		key := int(semi.Key())
		octave := uint32(semi.Octave())

		if key >= len(semitonePeriodTable) {
			return result
		}

		if c4SampleRate == 0 {
			c4SampleRate = period.Frequency(DefaultC4SampleRate)
		}

		if ft != 0 {
			c4SampleRate = CalcFinetuneC4SampleRate[period.Amiga](AmigaConverter, c4SampleRate, ft)
		}

		*p = period.Amiga(float64(floatDefaultC4SampleRate*semitonePeriodTable[key]) / float64(uint32(c4SampleRate)<<octave))
	default:
	}

	return result
}

// CalcFinetuneC4SampleRate calculates a new c4 sample rate after a finetune adjustment
func CalcFinetuneC4SampleRate[TPeriod period.Period](converter period.PeriodConverter[TPeriod], c4SampleRate period.Frequency, finetune note.Finetune) period.Frequency {
	if finetune == 0 {
		return c4SampleRate
	}

	nft := C4SlideFines + int(finetune)
	p := CalcSemitonePeriod[TPeriod](note.Semitone(nft/SlideFinesPerNote), note.Finetune(nft%SlideFinesPerNote), c4SampleRate)
	return converter.GetFrequency(p)
}

// FrequencyFromSemitone returns the frequency from the semitone (and c4 sample rate)
func FrequencyFromSemitone[TPeriod period.Period](converter period.PeriodConverter[TPeriod], semitone note.Semitone, c4SampleRate period.Frequency) float32 {
	p := CalcSemitonePeriod[TPeriod](semitone, 0, c4SampleRate)
	return float32(converter.GetFrequency(p))
}
