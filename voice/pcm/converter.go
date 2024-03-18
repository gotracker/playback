package pcm

import "github.com/gotracker/playback/mixing/volume"

// SampleConverter is an interface to a sample converter
type SampleConverter interface {
	Size() int
	ReadAt(s *SampleData, ofs int64) (volume.Volume, error)
	Format(s *SampleData) SampleDataFormat
}
