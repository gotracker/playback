package period

import (
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

const (
	// MiddleCFrequency is the default C2SPD for IT samples
	MiddleCFrequency = 8363
	// MiddleCPeriod is the sampler (Amiga-style) period of the C-5 note
	MiddleCPeriod = 856

	floatDefaultC2Spd = float32(MiddleCFrequency)

	// ITBaseClock is the base clock speed of IT files
	ITBaseClock period.Frequency = MiddleCFrequency * MiddleCPeriod

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
	if linearFreqSlides {
		nft := int(semi)*semitonesPerNote + int(ft)
		return Linear{
			// NOTE: not sure why the magic downshift a whole octave,
			// but it makes all the calculations work, so here we are.
			Finetune: note.Finetune(nft),
			C2Spd:    c2spd,
		}
	}

	key := int(semi.Key())
	octave := uint32(semi.Octave())

	if key >= len(semitonePeriodTable) {
		return nil
	}

	if c2spd == 0 {
		c2spd = period.Frequency(MiddleCFrequency)
	}

	if ft != 0 {
		c2spd = CalcFinetuneC2Spd(c2spd, ft, linearFreqSlides)
	}

	p := (Amiga(floatDefaultC2Spd*semitonePeriodTable[key]) / Amiga(uint32(c2spd)<<octave))
	p = p.AddInteger(0)
	return p
}

// CalcFinetuneC2Spd calculates a new C2SPD after a finetune adjustment
func CalcFinetuneC2Spd(c2spd period.Frequency, finetune note.Finetune, linearFreqSlides bool) period.Frequency {
	if finetune == 0 {
		return c2spd
	}

	nft := 5*semitonesPerOctave + int(finetune)
	p := CalcSemitonePeriod(note.Semitone(nft/semitonesPerNote), note.Finetune(nft%semitonesPerNote), c2spd, linearFreqSlides)
	return period.Frequency(p.GetFrequency())
}

// FrequencyFromSemitone returns the frequency from the semitone (and c2spd)
func FrequencyFromSemitone(semitone note.Semitone, c2spd period.Frequency, linearFreqSlides bool) float32 {
	p := CalcSemitonePeriod(semitone, 0, c2spd, linearFreqSlides)
	return float32(p.GetFrequency())
}
