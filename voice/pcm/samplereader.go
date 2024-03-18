package pcm

import "github.com/gotracker/playback/mixing/volume"

type SampleType interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~float32 | ~float64
}

type PCMReader[TConverter SampleConverter] struct {
	SampleData
	cnv TConverter
}

func (s *PCMReader[TConverter]) Read() (volume.Matrix, error) {
	return s.readData(s.cnv)
}

func (s PCMReader[TConverter]) Format() SampleDataFormat {
	return s.cnv.Format(&s.SampleData)
}
