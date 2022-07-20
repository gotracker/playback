package pcm

import "github.com/gotracker/gomixing/volume"

type SampleType interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~float32 | ~float64
}

type PCMReader[TConverter SampleConverter] struct {
	SampleData
	cnv TConverter
}

func (s *PCMReader[T]) Read() (volume.Matrix, error) {
	return s.readData(s.cnv)
}
