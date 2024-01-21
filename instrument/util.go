package instrument

import (
	"fmt"

	"github.com/gotracker/gomixing/sampling"
)

// ID is an identifier for an instrument/sample that means something to the format
type ID interface {
	IsEmpty() bool
	GetIndexAndSample() (int, int)
	fmt.Stringer
}

// Data is the interface to implementation-specific functions on an instrument
type Data interface {
	GetLength() sampling.Pos
}
