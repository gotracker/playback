package sampling

import (
	"encoding/binary"
	"io"

	"github.com/gotracker/playback/mixing/volume"
)

type Formatter interface {
	Size() int
	ReadAt(data []byte, ofs int64) (volume.Volume, error)
	WriteAt(data []byte, ofs int64, v volume.Volume) error
	Write(out io.Writer, v volume.Volume) error
}

func GetFormatter(format Format) Formatter {
	switch format {
	default:
		return nil
	case Format8BitUnsigned:
		// Format8BitUnsigned is for unsigned 8-bit data
		return Sample8BitUnsigned{}
	case Format8BitSigned:
		// Format8BitSigned is for signed 8-bit data
		return Sample8BitSigned{}
	case Format16BitLEUnsigned:
		// Format16BitLEUnsigned is for unsigned, little-endian, 16-bit data
		return Sample16BitUnsigned{byteOrder: binary.LittleEndian}
	case Format16BitLESigned:
		// Format16BitLESigned is for signed, little-endian, 16-bit data
		return Sample16BitSigned{byteOrder: binary.LittleEndian}
	case Format16BitBEUnsigned:
		// Format16BitBEUnsigned is for unsigned, big-endian, 16-bit data
		return Sample16BitUnsigned{byteOrder: binary.BigEndian}
	case Format16BitBESigned:
		// Format16BitBESigned is for signed, big-endian, 16-bit data
		return Sample16BitSigned{byteOrder: binary.BigEndian}
	case Format32BitLEFloat:
		// Format32BitLEFloat is for little-endian, 32-bit floating-point data
		return Sample32BitFloat{byteOrder: binary.LittleEndian}
	case Format32BitBEFloat:
		// Format32BitBEFloat is for big-endian, 32-bit floating-point data
		return Sample32BitFloat{byteOrder: binary.BigEndian}
	case Format64BitLEFloat:
		// Format64BitLEFloat is for little-endian, 64-bit floating-point data
		return Sample64BitFloat{byteOrder: binary.LittleEndian}
	case Format64BitBEFloat:
		// Format64BitBEFloat is for big-endian, 64-bit floating-point data
		return Sample64BitFloat{byteOrder: binary.BigEndian}
	}
}
