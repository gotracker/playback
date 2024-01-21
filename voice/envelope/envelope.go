package envelope

import (
	"github.com/gotracker/playback/voice/loop"
)

// Envelope is an envelope for instruments
type Envelope[T any] struct {
	Enabled bool
	Loop    loop.Loop
	Sustain loop.Loop
	Length  int
	Values  []Point[T]
}
