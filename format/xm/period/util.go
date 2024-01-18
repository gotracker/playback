package period

import (
	"math"

	xmSystem "github.com/gotracker/playback/format/xm/system"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
)

// CalcFinetuneC4SampleRate calculates a new c4 sample rate after a finetune adjustment
func CalcFinetuneC4SampleRate(c4SampleRate frequency.Frequency, st note.Semitone, finetune note.Finetune) frequency.Frequency {
	if finetune == 0 && st == xmSystem.C4Note {
		return c4SampleRate
	}

	per := max(float64(st)*xmSystem.SlideFinesPerNote+float64(finetune)/2, 0)
	exp := per / xmSystem.SlideFinesPerOctave
	pow := math.Pow(2.0, exp-xmSystem.C4Octave)

	freq := math.Floor(float64(c4SampleRate) * pow)

	return frequency.Frequency(freq)
}
