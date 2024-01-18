package system

import (
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/system"
)

const (
	// DefaultC4SampleRate is the default c4 sample rate for XM samples
	DefaultC4SampleRate = 8363
	// C4Period is the sampler (Amiga-style) period of the C-4 note
	C4Period = 856

	C4Octave = 4
	C4Note   = C4Octave * NotesPerOctave

	floatDefaultC4SampleRate = float32(DefaultC4SampleRate)

	// XMBaseClock is the base clock speed of XM files
	XMBaseClock frequency.Frequency = DefaultC4SampleRate * C4Period

	NotesPerOctave        = 12
	SlideFinesPerSemitone = 4
	SemitonesPerNote      = 16
	SlideFinesPerNote     = SlideFinesPerSemitone * SemitonesPerNote
	SlideFinesPerOctave   = SlideFinesPerNote * NotesPerOctave
	C4SlideFines          = C4Note * SlideFinesPerNote
)

var semitonePeriodTable = [...]uint16{27392, 25856, 24384, 23040, 21696, 20480, 19328, 18240, 17216, 16256, 15360, 14496}

var XMSystem system.ClockableSystem = system.ClockedSystem{
	MaxPastNotesPerChannel: 0,
	BaseClock:              XMBaseClock,
	BaseFinetunes:          C4SlideFines,
	FinetunesPerOctave:     SlideFinesPerOctave,
	FinetunesPerNote:       SlideFinesPerNote,
	CommonPeriod:           C4Period,
	CommonRate:             DefaultC4SampleRate,
	SemitonePeriods:        semitonePeriodTable,
}
