package channel

import (
	"fmt"
)

// InstID is an instrument ID in S3M world
type InstID uint8

// IsEmpty returns true if the instrument ID is 'nothing'
func (s InstID) IsEmpty() bool {
	return s == 0
}

func (s InstID) GetIndexAndSample() (int, int) {
	idx := int(s) - 1
	return idx, idx
}

func (s InstID) String() string {
	return fmt.Sprint(uint8(s))
}
