package instrument

import (
	"fmt"

	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/playback/note"
)

// ID is an identifier for an instrument/sample that means something to the format
type ID interface {
	IsEmpty() bool
	GetIndexAndSemitone() (int, note.Semitone)
	fmt.Stringer
}

// Data is the interface to implementation-specific functions on an instrument
type Data interface {
	GetLength() sampling.Pos
}
