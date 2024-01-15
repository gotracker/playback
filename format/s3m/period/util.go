package period

import (
	"math"

	"github.com/gotracker/playback/format/s3m/system"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

var DefaultC4SampleRate = system.DefaultC4SampleRate
var semitonePeriodTable = [...]float32{27392, 25856, 24384, 23040, 21696, 20480, 19328, 18240, 17216, 16256, 15360, 14496}

// CalcSemitonePeriod calculates the semitone period for it notes
func CalcSemitonePeriod(semi note.Semitone, ft note.Finetune, c4SampleRate frequency.Frequency) period.Amiga {
	if semi == note.UnchangedSemitone {
		panic("how?")
	}

	key := int(semi.Key())
	octave := uint32(semi.Octave())

	if key >= len(semitonePeriodTable) {
		var empty period.Amiga
		return empty
	}

	if c4SampleRate == 0 {
		c4SampleRate = frequency.Frequency(system.DefaultC4SampleRate)
	}

	if ft != 0 {
		c4SampleRate = CalcFinetuneC4SampleRate(c4SampleRate, ft)
	}

	return period.Amiga(float64(semitonePeriodTable[key]*float32(system.DefaultC4SampleRate)) / float64(uint32(c4SampleRate)<<octave))
}

// CalcFinetuneC4SampleRate calculates a new frequency after a finetune adjustment
func CalcFinetuneC4SampleRate(c4SampleRate frequency.Frequency, finetune note.Finetune) frequency.Frequency {
	if finetune == 0 {
		return c4SampleRate
	}

	nft := system.C4SlideFines + int(finetune)
	p := CalcSemitonePeriod(note.Semitone(nft/system.SlideFinesPerNote), note.Finetune(nft%system.SlideFinesPerNote), c4SampleRate)
	return p.GetFrequency()
}

// FrequencyFromSemitone returns the frequency from the semitone (and c4 sample rate)
func FrequencyFromSemitone(semitone note.Semitone, c4SampleRate frequency.Frequency) float32 {
	p := CalcSemitonePeriod(semitone, 0, c4SampleRate)
	return float32(p.GetFrequency())
}

// ToAmigaPeriod calculates an amiga period for a linear finetune period
func ToAmigaPeriod(finetunes note.Finetune, c4SampleRate frequency.Frequency) period.Amiga {
	if finetunes < 0 {
		finetunes = 0
	}
	pow := math.Pow(2, float64(finetunes)/system.SlideFinesPerOctave)
	linFreq := float64(c4SampleRate) * pow / float64(system.DefaultC4SampleRate)

	return period.Amiga(float64(semitonePeriodTable[0]) / linFreq)
}
