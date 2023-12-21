package channel

import (
	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/playback/note"
)

// NoteFactory produces a note effect for the provided channel pattern data
func NoteFactory(mem *Memory, n s3mfile.Note, i InstID) EffectS3M {
	switch {
	case n == s3mfile.EmptyNote:
		if i != 0 {
			return start(0, i)
		}
		return nil
	case n == s3mfile.StopNote:
		return stop()
	default:
		k := uint8(n.Key()) & 0x0f
		o := uint8(n.Octave()) & 0x0f
		if k < 12 && o < 10 {
			s := note.Semitone(o*12 + k)
			return start(s, i)
		}
	}
	return note.InvalidNote{}
}
