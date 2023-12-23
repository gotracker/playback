package period

import (
	"math"

	"github.com/gotracker/playback/format/it/system"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

var DefaultC5SampleRate = system.DefaultC5SampleRate
var semitonePeriodTable = [...]float32{27392, 25856, 24384, 23040, 21696, 20480, 19328, 18240, 17216, 16256, 15360, 14496}

// CalcSemitonePeriod calculates the semitone period for it notes
func CalcSemitonePeriod[TPeriod period.Period](semi note.Semitone, ft note.Finetune, c5SampleRate period.Frequency) TPeriod {
	if semi == note.UnchangedSemitone {
		panic("how?")
	}
	var result TPeriod
	switch p := any(&result).(type) {
	case *period.Linear:
		nft := int(semi)*system.SlideFinesPerNote + int(ft)
		p.Finetune = note.Finetune(nft)
		p.CommonRate = c5SampleRate
	case *period.Amiga:
		key := int(semi.Key())
		octave := uint32(semi.Octave())

		if key >= len(semitonePeriodTable) {
			return result
		}

		if c5SampleRate == 0 {
			c5SampleRate = period.Frequency(system.DefaultC5SampleRate)
		}

		if ft != 0 {
			c5SampleRate = CalcFinetuneC5SampleRate[period.Amiga](AmigaConverter, c5SampleRate, ft)
		}

		*p = period.Amiga(float64(semitonePeriodTable[key]*system.DefaultC5SampleRate) / float64(uint32(c5SampleRate)<<octave))
	default:
	}

	return result
}

// CalcFinetuneC5SampleRate calculates a new frequency after a finetune adjustment
func CalcFinetuneC5SampleRate[TPeriod period.Period](converter period.PeriodConverter[TPeriod], c5SampleRate period.Frequency, finetune note.Finetune) period.Frequency {
	if finetune == 0 {
		return c5SampleRate
	}

	nft := system.C5SlideFines + int(finetune)
	p := CalcSemitonePeriod[TPeriod](note.Semitone(nft/system.SlideFinesPerNote), note.Finetune(nft%system.SlideFinesPerNote), c5SampleRate)
	return converter.GetFrequency(p)
}

// FrequencyFromSemitone returns the frequency from the semitone (and c5 sample rate)
func FrequencyFromSemitone[TPeriod period.Period](converter period.PeriodConverter[TPeriod], semitone note.Semitone, c5SampleRate period.Frequency) float32 {
	p := CalcSemitonePeriod[TPeriod](semitone, 0, c5SampleRate)
	return float32(converter.GetFrequency(p))
}

// ToAmigaPeriod calculates an amiga period for a linear finetune period
func ToAmigaPeriod(finetunes note.Finetune, c5SampleRate period.Frequency) period.Amiga {
	if finetunes < 0 {
		finetunes = 0
	}
	pow := math.Pow(2, float64(finetunes)/system.SlideFinesPerOctave)
	linFreq := float64(c5SampleRate) * pow / float64(system.DefaultC5SampleRate)

	return period.Amiga(float64(semitonePeriodTable[0]) / linFreq)
}
