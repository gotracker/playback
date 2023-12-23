package system

import (
	"github.com/gotracker/playback/system"
)

const (
	// DefaultC5SampleRate is the default sample rate for IT samples
	DefaultC5SampleRate = 8363
	// C5Period is the sampler (Amiga-style) period of the C-5 note
	C5Period = 428

	// ITBaseClock is the base clock speed of IT files
	ITBaseClock system.Frequency = DefaultC5SampleRate * C5Period

	NotesPerOctave        = 12
	SlideFinesPerSemitone = 4
	SemitonesPerNote      = 16
	SlideFinesPerNote     = SlideFinesPerSemitone * SemitonesPerNote
	SlideFinesPerOctave   = SlideFinesPerNote * NotesPerOctave
	C5SlideFines          = 5 * SlideFinesPerOctave
)

var ITSystem system.System = system.ClockedSystem{
	BaseClock:          ITBaseClock,
	BaseFinetunes:      C5SlideFines,
	FinetunesPerOctave: SlideFinesPerOctave,
}
