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

var semitonePeriodTable = [...]float32{27392, 25856, 24384, 23040, 21696, 20480, 19328, 18240, 17216, 16256, 15360, 14496}

var S3MSystem system.ClockableSystem = system.ClockedSystem{
	MaxPastNotesPerChannel: 0,
	BaseClock:              S3MBaseClock,
	BaseFinetunes:          C4SlideFines,
	FinetunesPerOctave:     SlideFinesPerOctave,
	FinetunesPerNote:       SlideFinesPerNote,
	CommonRate:             DefaultC4SampleRate,
	SemitonePeriods:        semitonePeriodTable,
}
