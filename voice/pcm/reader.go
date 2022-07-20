package pcm

import (
	"github.com/gotracker/gomixing/volume"
)

// SampleReader is a reader interface that can return a whole multichannel sample at the current position
type SampleReader interface {
	Read() (volume.Matrix, error)
}

func (s *SampleData) readData(converter SampleConverter) (volume.Matrix, error) {
	bps := converter.Size()
	actualPos := int64(s.pos * s.channels * bps)
	if actualPos < 0 {
		actualPos = 0
	}

	out := volume.Matrix{
		Channels: s.channels,
	}
	for c := 0; c < s.channels; c++ {
		v, err := converter.ReadAt(s, actualPos)
		if err != nil {
			return volume.Matrix{}, err
		}

		out.StaticMatrix[c] = v
		actualPos += int64(bps)
	}

	s.pos++
	return out, nil
}
