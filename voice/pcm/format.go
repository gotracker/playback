package pcm

import "math"

// SampleDataFormat is the format of the sample data
type SampleDataFormat uint8

const (
	// SampleDataFormat8BitUnsigned is for unsigned 8-bit data
	SampleDataFormat8BitUnsigned = SampleDataFormat(iota)
	// SampleDataFormat8BitSigned is for signed 8-bit data
	SampleDataFormat8BitSigned
	// SampleDataFormat16BitLEUnsigned is for unsigned, little-endian, 16-bit data
	SampleDataFormat16BitLEUnsigned
	// SampleDataFormat16BitLESigned is for signed, little-endian, 16-bit data
	SampleDataFormat16BitLESigned
	// SampleDataFormat16BitBEUnsigned is for unsigned, big-endian, 16-bit data
	SampleDataFormat16BitBEUnsigned
	// SampleDataFormat16BitBESigned is for signed, big-endian, 16-bit data
	SampleDataFormat16BitBESigned
	// SampleDataFormat32BitLEFloat is for little-endian, 32-bit floating-point data
	SampleDataFormat32BitLEFloat
	// SampleDataFormat32BitBEFloat is for big-endian, 32-bit floating-point data
	SampleDataFormat32BitBEFloat
	// SampleDataFormat64BitLEFloat is for little-endian, 64-bit floating-point data
	SampleDataFormat64BitLEFloat
	// SampleDataFormat64BitBEFloat is for big-endian, 64-bit floating-point data
	SampleDataFormat64BitBEFloat
)

const SampleDataFormatNative = math.MaxUint8

func getSampleBytes(sdf SampleDataFormat) int {
	switch sdf {
	case SampleDataFormat8BitUnsigned:
		return Sample8BitUnsigned{}.Size()

	case SampleDataFormat8BitSigned:
		return Sample8BitSigned{}.Size()

	case SampleDataFormat16BitLEUnsigned, SampleDataFormat16BitBEUnsigned:
		return Sample16BitUnsigned{}.Size()

	case SampleDataFormat16BitLESigned, SampleDataFormat16BitBESigned:
		return Sample16BitSigned{}.Size()

	case SampleDataFormat32BitLEFloat, SampleDataFormat32BitBEFloat:
		return Sample32BitFloat{}.Size()

	case SampleDataFormat64BitLEFloat, SampleDataFormat64BitBEFloat:
		return Sample64BitFloat{}.Size()

	default:
		return 1
	}
}
