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
	Values  []EnvPoint[T]
}

// EnvPoint is a point for the envelope
type EnvPoint[T any] struct {
	Pos    int
	Length int
	Y      T
}
