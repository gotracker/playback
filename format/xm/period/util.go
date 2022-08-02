package period

import (
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

const (
	// MiddleCFrequency is the default C2SPD for XM samples
	MiddleCFrequency = 8363
	// MiddleCPeriod is the sampler (Amiga-style) period of the C-5 note
	MiddleCPeriod = 1712

	// BaseClock is the base clock speed of XM files
	BaseClock period.Frequency = MiddleCFrequency * MiddleCPeriod

	notesPerOctave     = 12
	semitonesPerNote   = 64
	semitonesPerOctave = notesPerOctave * semitonesPerNote
)

var semitonePeriodTable = [...]float32{27392, 25856, 24384, 23040, 21696, 20480, 19328, 18240, 17216, 16256, 15360, 14496}

// CalcSemitonePeriod calculates the semitone period for it notes
func CalcSemitonePeriod(semi note.Semitone, ft note.Finetune, c2spd period.Frequency, linearFreqSlides bool) period.Period {
	if semi == note.UnchangedSemitone {
		panic("how?")
	}

	if c2spd == 0 {
		c2spd = period.Frequency(MiddleCFrequency)
	}

	nft := note.Finetune(semi)*semitonesPerNote + ft
	if linearFreqSlides {
		return Linear{
			Finetune: nft,
			C2Spd:    c2spd,
		}
	}

	return ToAmigaPeriod(nft, c2spd).AddInteger(0)
}

// CalcFinetuneC2Spd calculates a new C2SPD after a finetune adjustment
func CalcFinetuneC2Spd(c2spd period.Frequency, finetune note.Finetune, linearFreqSlides bool) period.Frequency {
	if finetune == 0 {
		return c2spd
	}

	nft := 4*semitonesPerOctave + int(finetune)
	p := CalcSemitonePeriod(note.Semitone(nft/semitonesPerNote), note.Finetune(nft%semitonesPerNote), c2spd, linearFreqSlides)
	return p.GetFrequency()
}

// FrequencyFromSemitone returns the frequency from the semitone (and c2spd)
func FrequencyFromSemitone(semitone note.Semitone, c2spd period.Frequency, linearFreqSlides bool) float32 {
	p := CalcSemitonePeriod(semitone, 0, c2spd, linearFreqSlides)
	return float32(p.GetFrequency())
}
