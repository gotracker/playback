package sampling

import "github.com/gotracker/playback/mixing/volume"

// Sampler is an interface to the sampling system
type Sampler interface {
	GetPosition() Pos
	Advance()
	GetSample() volume.Matrix
}

// SampleStream is an interface to a sample stream (aka: an instrument)
type SampleStream interface {
	GetSample(Pos) volume.Matrix
}

// NewSampler creates a basic sampler that implements the Sampler interface
func NewSampler(ss SampleStream, pos Pos, period float32) Sampler {
	s := sampler{
		ss:     ss,
		pos:    pos,
		period: period,
	}
	return &s
}
