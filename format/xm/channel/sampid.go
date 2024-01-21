package channel

import (
	"fmt"
)

// SampleID is an InstrumentID that is a combination of InstID and SampID
type SampleID struct {
	InstID uint8
	SampID uint8
}

// IsEmpty returns true if the sample ID is empty
func (s SampleID) IsEmpty() bool {
	return s.InstID == 0
}

func (s SampleID) GetIndexAndSample() (int, int) {
	return int(s.InstID) - 1, int(s.SampID)
}

func (s SampleID) String() string {
	return fmt.Sprintf("%d(%d)", s.InstID, s.SampID)
}
