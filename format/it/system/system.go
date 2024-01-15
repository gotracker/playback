package system

import (
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/system"
)

const (
	// DefaultC5SampleRate is the default sample rate for IT samples
	DefaultC5SampleRate = 8363
	// C5Period is the sampler (Amiga-style) period of the C-5 note
	C5Period = 428

	// ITBaseClock is the base clock speed of IT files
	ITBaseClock frequency.Frequency = DefaultC5SampleRate * C5Period

	NotesPerOctave        = 12
	SlideFinesPerSemitone = 4
	SemitonesPerNote      = 16
	SlideFinesPerNote     = SlideFinesPerSemitone * SemitonesPerNote
	SlideFinesPerOctave   = SlideFinesPerNote * NotesPerOctave
	C5SlideFines          = 5 * SlideFinesPerOctave
)

var semitonePeriodTable = [...]float32{27392, 25856, 24384, 23040, 21696, 20480, 19328, 18240, 17216, 16256, 15360, 14496}

var ITSystem system.ClockableSystem = system.ClockedSystem{
	MaxPastNotesPerChannel: 1,
	BaseClock:              ITBaseClock,
	BaseFinetunes:          C5SlideFines,
	FinetunesPerOctave:     SlideFinesPerOctave,
	FinetunesPerNote:       SlideFinesPerNote,
	CommonRate:             DefaultC5SampleRate,
	SemitonePeriods:        semitonePeriodTable,
}
