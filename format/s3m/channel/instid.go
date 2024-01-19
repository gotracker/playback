package channel

import (
	"fmt"

	"github.com/gotracker/playback/note"
)

// InstID is an instrument ID in S3M world
type InstID uint8

// IsEmpty returns true if the instrument ID is 'nothing'
func (s InstID) IsEmpty() bool {
	return s == 0
}

func (s InstID) GetIndexAndSemitone() (int, note.Semitone) {
	return int(s) - 1, note.UnchangedSemitone
}

func (s InstID) String() string {
	return fmt.Sprint(uint8(s))
}
