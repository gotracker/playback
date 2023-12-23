package system

import (
	"github.com/gotracker/playback/system"
)

const (
	// DefaultC4SampleRate is the default c4 sample rate for XM samples
	DefaultC4SampleRate = 8363
	// C4Period is the sampler (Amiga-style) period of the C-4 note
	C4Period = 1712

	floatDefaultC4SampleRate = float32(DefaultC4SampleRate)

	// XMBaseClock is the base clock speed of XM files
	XMBaseClock system.Frequency = DefaultC4SampleRate * C4Period

	NotesPerOctave        = 12
	SlideFinesPerSemitone = 4
	SemitonesPerNote      = 16
	SlideFinesPerNote     = SlideFinesPerSemitone * SemitonesPerNote
	SlideFinesPerOctave   = SlideFinesPerNote * NotesPerOctave
	C4SlideFines          = 4 * SlideFinesPerOctave
)

var XMSystem system.System = system.ClockedSystem{
	BaseClock:          XMBaseClock,
	BaseFinetunes:      C4SlideFines,
	FinetunesPerOctave: SlideFinesPerOctave,
}
