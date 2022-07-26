package channel

import (
	"fmt"

	"github.com/gotracker/playback/note"
)

// SampleID is an InstrumentID that is a combination of InstID and SampID
type SampleID struct {
	InstID   uint8
	Semitone note.Semitone
}

// IsEmpty returns true if the sample ID is empty
func (s SampleID) IsEmpty() bool {
	return s.InstID == 0
}

func (s SampleID) String() string {
	return fmt.Sprint(s.InstID)
}

type SemitoneAndSampleID struct {
	ST note.Semitone
	ID uint8
}
