package frequency

import (
	"fmt"
)

// Frequency is a frequency value, in Hertz (Hz)
type Frequency float64

func (f Frequency) GoString() string {
	switch {
	case f < 1_000:
		return fmt.Sprintf("%fHz", f)
	case f < 1_000_000:
		return fmt.Sprintf("%3.fkHz", f)
	case f < 1_000_000_000:
		return fmt.Sprintf("%3.fMHz", f)
	case f < 1_000_000_000_000:
		return fmt.Sprintf("%3.fGHz", f)
	case f < 1_000_000_000_000_000:
		return fmt.Sprintf("%3.fTHz", f)
	default:
		return fmt.Sprintf("%fHz", f)
	}
}
