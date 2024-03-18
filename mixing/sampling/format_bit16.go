package sampling

import (
	"encoding/binary"
	"io"

	"github.com/gotracker/playback/mixing/volume"
)

const (
	cSample16BitDataCoeff   = 0x8000
	cSample16BitVolumeCoeff = volume.Volume(1) / cSample16BitDataCoeff
	cSample16BitBytes       = 2
)

// Sample16BitSigned is a signed 16-bit sample
type Sample16BitSigned struct {
	byteOrder binary.ByteOrder
}

// Volume returns the volume value for the sample
func (Sample16BitSigned) volume(v int16) volume.Volume {
	return volume.Volume(v) * cSample16BitVolumeCoeff
}

// fromVolume returns the volume value for the sample
func (Sample16BitSigned) fromVolume(v volume.Volume) int16 {
	return int16(v.ToIntSample(16))
}

// Size returns the size of the sample in bytes
func (Sample16BitSigned) Size() int {
	return cSample16BitBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s Sample16BitSigned) ReadAt(data []byte, ofs int64) (volume.Volume, error) {
	if len(data) <= int(ofs)+(cSample16BitBytes-1) {
		return 0, io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	v := int16(s.byteOrder.Uint16(data[ofs:]))
	return s.volume(v), nil
}

// WriteAt writes a value to the slice provided in the byte order provided
func (s Sample16BitSigned) WriteAt(data []byte, ofs int64, v volume.Volume) error {
	if len(data) <= int(ofs) {
		return io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	s.byteOrder.PutUint16(data[ofs:], uint16(s.fromVolume(v)))
	return nil
}

// Write writes a value to the Writer provided in the byte order provided
func (s Sample16BitSigned) Write(out io.Writer, v volume.Volume) error {
	return binary.Write(out, s.byteOrder, s.fromVolume(v))
}

// Sample16BitUnsigned is an unsigned 16-bit sample
type Sample16BitUnsigned struct {
	byteOrder binary.ByteOrder
}

// Volume returns the volume value for the sample
func (Sample16BitUnsigned) volume(v uint16) volume.Volume {
	return volume.Volume(int16(v-cSample16BitDataCoeff)) * cSample16BitVolumeCoeff
}

// fromVolume returns the volume value for the sample
func (Sample16BitUnsigned) fromVolume(v volume.Volume) uint16 {
	return uint16(v.ToUintSample(16))
}

// Size returns the size of the sample in bytes
func (Sample16BitUnsigned) Size() int {
	return cSample16BitBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s Sample16BitUnsigned) ReadAt(data []byte, ofs int64) (volume.Volume, error) {
	if len(data) <= int(ofs)+(cSample16BitBytes-1) {
		return 0, io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	v := uint16(s.byteOrder.Uint16(data[ofs:]))
	return s.volume(v), nil
}

// WriteAt writes a value to the slice provided in the byte order provided
func (s Sample16BitUnsigned) WriteAt(data []byte, ofs int64, v volume.Volume) error {
	if len(data) <= int(ofs) {
		return io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	s.byteOrder.PutUint16(data[ofs:], s.fromVolume(v))
	return nil
}

// Write writes a value to the Writer provided in the byte order provided
func (s Sample16BitUnsigned) Write(out io.Writer, v volume.Volume) error {
	return binary.Write(out, s.byteOrder, s.fromVolume(v))
}
