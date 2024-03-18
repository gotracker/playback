package pcm

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"math"

	"github.com/gotracker/playback/mixing/volume"
)

var (
	// ErrIndexOutOfRange is for when a slice is iterated with an index that's out of the range
	ErrIndexOutOfRange = errors.New("index out of range")
)

type NativeSampleData struct {
	baseSampleData
	data []volume.Matrix
}

// SampleReaderNative is a native (pre-converted) PCM sample reader
type SampleReaderNative struct {
	NativeSampleData
}

// Channels returns the channel count from the sample data
func (s *NativeSampleData) Channels() int {
	return s.channels
}

// Length returns the sample length from the sample data
func (s *NativeSampleData) Length() int {
	return s.length
}

// Seek sets the current position in the sample data
func (s *NativeSampleData) Seek(pos int) {
	s.pos = pos
}

// Tell returns the current position in the sample data
func (s *NativeSampleData) Tell() int {
	return s.pos
}

// Read returns the next multichannel sample
func (s *SampleReaderNative) Read() (volume.Matrix, error) {
	return s.readData()
}

func (s *NativeSampleData) readData() (volume.Matrix, error) {
	if s.pos < 0 {
		s.pos = 0
	}

	if s.pos >= s.length {
		return volume.Matrix{}, ErrIndexOutOfRange
	}
	samp := s.data[s.pos]
	s.pos++

	return samp, nil
}

func (s NativeSampleData) Format() SampleDataFormat {
	return SampleDataFormat64BitLEFloat
}

func (s NativeSampleData) Base64() string {
	var src bytes.Buffer
	for _, d := range s.data {
		for c := 0; c < s.channels; c++ {
			binary.Write(&src, binary.LittleEndian, math.Float64bits(float64(d.StaticMatrix[c])))
		}
	}

	return base64.StdEncoding.EncodeToString(src.Bytes())
}

func NewSampleNative(data []volume.Matrix, length int, channels int) Sample {
	return &SampleReaderNative{
		NativeSampleData: NativeSampleData{
			baseSampleData: baseSampleData{
				length:   length,
				channels: channels,
			},
			data: data,
		},
	}
}
