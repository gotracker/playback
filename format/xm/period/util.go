package period

import (
	"math"

	xmSystem "github.com/gotracker/playback/format/xm/system"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/system"
)

// CalcSemitonePeriod calculates the semitone period for it notes
func CalcSemitonePeriod[TPeriod period.Period](semi note.Semitone, ft note.Finetune, c4SampleRate period.Frequency) TPeriod {
	if semi == note.UnchangedSemitone {
		panic("how?")
	}

	var result TPeriod
	switch p := any(&result).(type) {
	case *period.Linear:
		nft := int(semi)*xmSystem.SlideFinesPerNote + int(ft)/2
		p.Finetune = note.Finetune(nft)
	case *period.Amiga:
		stp, ok := xmSystem.XMSystem.GetSemitonePeriod(semi.Key())
		if !ok {
			return result
		}

		if c4SampleRate == 0 {
			c4SampleRate = period.Frequency(xmSystem.DefaultC4SampleRate)
		}

		per := max(float64(xmSystem.C4Note)*xmSystem.SlideFinesPerNote-float64(ft)/2, 0)

		exp := (per / xmSystem.SlideFinesPerOctave) - xmSystem.C4Octave + float64(semi.Octave())

		pow := math.Pow(2.0, exp)
		sampleRate := math.Floor(float64(c4SampleRate) * pow)

		if sampleRate == 0 {
			return result
		}

		const defaultC4SampleRate = float64(xmSystem.DefaultC4SampleRate)

		ratio := defaultC4SampleRate / float64(sampleRate)

		*p = period.Amiga(ratio * float64(stp))
	default:
	}

	return result
}

// CalcFinetuneC4SampleRate calculates a new c4 sample rate after a finetune adjustment
func CalcFinetuneC4SampleRate(c4SampleRate period.Frequency, st note.Semitone, finetune note.Finetune) period.Frequency {
	if finetune == 0 && st == xmSystem.C4Note {
		return c4SampleRate
	}

	per := max(float64(st)*xmSystem.SlideFinesPerNote+float64(finetune)/2, 0)
	exp := per / xmSystem.SlideFinesPerOctave
	pow := math.Pow(2.0, exp-xmSystem.C4Octave)

	freq := math.Floor(float64(c4SampleRate) * pow)

	return system.Frequency(freq)
}

// FrequencyFromSemitone returns the frequency from the semitone (and c4 sample rate)
func FrequencyFromSemitone[TPeriod period.Period](converter period.PeriodConverter[TPeriod], semitone note.Semitone, c4SampleRate period.Frequency) float32 {
	p := CalcSemitonePeriod[TPeriod](semitone, 0, c4SampleRate)
	return float32(converter.GetFrequency(p))
}
