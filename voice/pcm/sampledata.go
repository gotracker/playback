package pcm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"

	"github.com/gotracker/gomixing/volume"
)

// Sample is the interface to a sample
type Sample interface {
	SampleReader
	Channels() int
	Length() int
	Seek(pos int)
	Tell() int
}

// SampleData is the presentation of the core data of the sample
type baseSampleData struct {
	pos      int // in multichannel samples
	length   int // in multichannel samples
	channels int
}

type SampleData struct {
	baseSampleData
	byteOrder binary.ByteOrder
	data      []byte
}

// Channels returns the channel count from the sample data
func (s *SampleData) Channels() int {
	return s.channels
}

// Length returns the sample length from the sample data
func (s *SampleData) Length() int {
	return s.length
}

// Seek sets the current position in the sample data
func (s *SampleData) Seek(pos int) {
	s.pos = pos
}

// Tell returns the current position in the sample data
func (s *SampleData) Tell() int {
	return s.pos
}

// NewSample constructs a sampler that can handle the requested sampler format
func NewSample(data []byte, length int, channels int, format SampleDataFormat) Sample {
	base := baseSampleData{
		length:   length,
		channels: channels,
	}
	switch format {
	case SampleDataFormat8BitSigned:
		return &PCMReader[Sample8BitSigned]{
			SampleData: SampleData{
				baseSampleData: base,
				byteOrder:      binary.LittleEndian,
				data:           data,
			},
		}
	case SampleDataFormat8BitUnsigned:
		return &PCMReader[Sample8BitUnsigned]{
			SampleData: SampleData{
				baseSampleData: base,
				byteOrder:      binary.LittleEndian,
				data:           data,
			},
		}
	case SampleDataFormat16BitLESigned:
		return &PCMReader[Sample16BitSigned]{
			SampleData: SampleData{
				baseSampleData: base,
				byteOrder:      binary.LittleEndian,
				data:           data,
			},
		}
	case SampleDataFormat16BitLEUnsigned:
		return &PCMReader[Sample16BitUnsigned]{
			SampleData: SampleData{
				baseSampleData: base,
				byteOrder:      binary.LittleEndian,
				data:           data,
			},
		}
	case SampleDataFormat16BitBESigned:
		return &PCMReader[Sample16BitSigned]{
			SampleData: SampleData{
				baseSampleData: base,
				byteOrder:      binary.BigEndian,
				data:           data,
			},
		}
	case SampleDataFormat16BitBEUnsigned:
		return &PCMReader[Sample16BitUnsigned]{
			SampleData: SampleData{
				baseSampleData: base,
				byteOrder:      binary.BigEndian,
				data:           data,
			},
		}
	case SampleDataFormat32BitLEFloat:
		return &PCMReader[Sample32BitFloat]{
			SampleData: SampleData{
				baseSampleData: base,
				byteOrder:      binary.LittleEndian,
				data:           data,
			},
		}
	case SampleDataFormat32BitBEFloat:
		return &PCMReader[Sample32BitFloat]{
			SampleData: SampleData{
				baseSampleData: base,
				byteOrder:      binary.BigEndian,
				data:           data,
			},
		}
	case SampleDataFormat64BitLEFloat:
		return &PCMReader[Sample64BitFloat]{
			SampleData: SampleData{
				baseSampleData: base,
				byteOrder:      binary.LittleEndian,
				data:           data,
			},
		}
	case SampleDataFormat64BitBEFloat:
		return &PCMReader[Sample64BitFloat]{
			SampleData: SampleData{
				baseSampleData: base,
				byteOrder:      binary.BigEndian,
				data:           data,
			},
		}
	default:
		panic("unhandled sampler type")
	}
}

func ConvertTo(from Sample, format SampleDataFormat) (Sample, error) {
	cvt := &bytes.Buffer{}
	length := from.Length()
	channels := from.Channels()
	for i := 0; i < length; i++ {
		samp, _ := from.Read() // ignore error
		for c := 0; c < channels; c++ {
			var vol volume.Volume
			if samp.Channels > c {
				vol = samp.StaticMatrix[c]
			}
			switch format {
			case SampleDataFormat8BitUnsigned:
				cv := (vol * 0x80) + 0x80
				if err := binary.Write(cvt, binary.LittleEndian, uint8(cv)); err != nil {
					return nil, err
				}
			case SampleDataFormat8BitSigned:
				cv := (vol * 0x80)
				if err := binary.Write(cvt, binary.LittleEndian, int8(cv)); err != nil {
					return nil, err
				}
			case SampleDataFormat16BitLEUnsigned:
				cv := (vol * 0x8000) + 0x8000
				if err := binary.Write(cvt, binary.LittleEndian, uint16(cv)); err != nil {
					return nil, err
				}
			case SampleDataFormat16BitLESigned:
				cv := (vol * 0x8000)
				if err := binary.Write(cvt, binary.LittleEndian, int16(cv)); err != nil {
					return nil, err
				}
			case SampleDataFormat16BitBEUnsigned:
				cv := (vol * 0x8000) + 0x8000
				if err := binary.Write(cvt, binary.BigEndian, uint16(cv)); err != nil {
					return nil, err
				}
			case SampleDataFormat16BitBESigned:
				cv := (vol * 0x8000)
				if err := binary.Write(cvt, binary.BigEndian, int16(cv)); err != nil {
					return nil, err
				}
			case SampleDataFormat32BitLEFloat:
				cv := vol
				if err := binary.Write(cvt, binary.LittleEndian, math.Float32bits(float32(cv))); err != nil {
					return nil, err
				}
			case SampleDataFormat32BitBEFloat:
				cv := vol
				if err := binary.Write(cvt, binary.BigEndian, math.Float32bits(float32(cv))); err != nil {
					return nil, err
				}
			case SampleDataFormat64BitLEFloat:
				cv := vol
				if err := binary.Write(cvt, binary.LittleEndian, math.Float64bits(float64(cv))); err != nil {
					return nil, err
				}
			case SampleDataFormat64BitBEFloat:
				cv := vol
				if err := binary.Write(cvt, binary.BigEndian, math.Float64bits(float64(cv))); err != nil {
					return nil, err
				}
			default:
				return nil, errors.New("unhandled format type")
			}
		}
	}
	to := NewSample(cvt.Bytes(), length, channels, format)
	return to, nil
}
