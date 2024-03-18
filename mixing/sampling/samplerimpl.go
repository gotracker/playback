package sampling

import "github.com/gotracker/playback/mixing/volume"

type sampler struct {
	ss     SampleStream
	pos    Pos
	period float32
}

func (s *sampler) GetPosition() Pos {
	return s.pos
}

func (s *sampler) Advance() {
	s.pos.Add(s.period)
}

func (s *sampler) GetSample() volume.Matrix {
	if s.ss == nil {
		return volume.Matrix{}
	}
	return s.ss.GetSample(s.pos)
}
