package system

import (
	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"

	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/system"
)

const (
	floatDefaultC4SampleRate = float32(DefaultC4SampleRate)
	c2Period                 = 1712

	// DefaultC4SampleRate is the default c4 sample rate for S3M samples
	DefaultC4SampleRate = system.Frequency(s3mfile.DefaultC2Spd)

	// S3MBaseClock is the base clock speed of S3M files
	S3MBaseClock period.Frequency = DefaultC4SampleRate * c2Period

	NotesPerOctave        = 12
	SlideFinesPerSemitone = 4
	SemitonesPerNote      = 16
	SlideFinesPerNote     = SlideFinesPerSemitone * SemitonesPerNote
	SlideFinesPerOctave   = SlideFinesPerNote * NotesPerOctave

	C4SlideFines = 4 * SlideFinesPerOctave
)

var S3MSystem system.System = system.ClockedSystem{
	BaseClock:          S3MBaseClock,
	BaseFinetunes:      C4SlideFines,
	FinetunesPerOctave: SlideFinesPerOctave,
}
