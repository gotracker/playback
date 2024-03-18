package sampling

import (
	"encoding/binary"
	"io"

	"github.com/gotracker/playback/mixing/volume"
)

const (
	cSample8BitDataCoeff   = 0x80
	cSample8BitVolumeCoeff = volume.Volume(1) / cSample8BitDataCoeff
	cSample8BitBytes       = 1
)

// Sample8BitSigned is a signed 8-bit sample
type Sample8BitSigned struct{}

// toVolume returns the volume value for the sample
func (Sample8BitSigned) toVolume(v int8) volume.Volume {
	return volume.Volume(v) * cSample8BitVolumeCoeff
}

// fromVolume returns the volume value for the sample
func (Sample8BitSigned) fromVolume(v volume.Volume) int8 {
	return int8(v.ToIntSample(8))
}

// Size returns the size of the sample in bytes
func (Sample8BitSigned) Size() int {
	return cSample8BitBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s Sample8BitSigned) ReadAt(data []byte, ofs int64) (volume.Volume, error) {
	if len(data) <= int(ofs) {
		return 0, io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	v := int8(data[ofs])
	return s.toVolume(v), nil
}

// WriteAt writes a value to the slice provided in the byte order provided
func (s Sample8BitSigned) WriteAt(data []byte, ofs int64, v volume.Volume) error {
	if len(data) <= int(ofs) {
		return io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	data[ofs] = uint8(s.fromVolume(v))
	return nil
}

// Write writes a value to the Writer provided in the byte order provided
func (s Sample8BitSigned) Write(out io.Writer, v volume.Volume) error {
	return binary.Write(out, binary.LittleEndian, s.fromVolume(v))
}

// Sample8BitUnsigned is an unsigned 8-bit sample
type Sample8BitUnsigned struct{}

// toVolume returns the volume value for the sample
func (Sample8BitUnsigned) toVolume(v uint8) volume.Volume {
	return volume.Volume(int8(v-uint8(cSample8BitDataCoeff))) * cSample8BitVolumeCoeff
}

// fromVolume returns the volume value for the sample
func (Sample8BitUnsigned) fromVolume(v volume.Volume) uint8 {
	return uint8(v.ToUintSample(8))
}

// Size returns the size of the sample in bytes
func (Sample8BitUnsigned) Size() int {
	return cSample8BitBytes
}

// ReadAt reads a value from the slice provided in the byte order provided
func (s Sample8BitUnsigned) ReadAt(data []byte, ofs int64) (volume.Volume, error) {
	if len(data) <= int(ofs) {
		return 0, io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	v := uint8(data[ofs])
	return s.toVolume(v), nil
}

// WriteAt writes a value to the slice provided in the byte order provided
func (s Sample8BitUnsigned) WriteAt(data []byte, ofs int64, v volume.Volume) error {
	if len(data) <= int(ofs) {
		return io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	data[ofs] = s.fromVolume(v)
	return nil
}

// Write writes a value to the Writer provided in the byte order provided
func (s Sample8BitUnsigned) Write(out io.Writer, v volume.Volume) error {
	return binary.Write(out, binary.LittleEndian, s.fromVolume(v))
}
