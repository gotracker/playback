package channel

import (
	"fmt"

	"github.com/gotracker/playback/note"
	"github.com/heucuva/optional"
)

// SampleID is an InstrumentID that is a combination of InstID and SampID
type SampleID struct {
	InstID   uint8
	Semitone optional.Value[note.Semitone]
}

// IsEmpty returns true if the sample ID is empty
func (s SampleID) IsEmpty() bool {
	return s.InstID == 0
}

func (s SampleID) GetIndexAndSemitone() (int, note.Semitone) {
	st := note.UnchangedSemitone
	if ost, set := s.Semitone.Get(); set {
		st = ost
	}
	return int(s.InstID) - 1, st
}

func (s SampleID) String() string {
	return fmt.Sprint(s.InstID)
}
