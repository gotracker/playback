package volume

import (
	"fmt"
	"math"
)

// Volume is a mixable volume
type Volume float32

const denormThreshold Volume = 1e-18

func clearDenorm(v Volume) Volume {
	if v > -denormThreshold && v < denormThreshold {
		return 0
	}
	return v
}

var (
	// VolumeUseInstVol tells the system to use the volume stored on the instrument
	// This is useful for trackers and other musical applications
	VolumeUseInstVol = Volume(math.Inf(-1))
)

// Int24 is an approximation of a 24-bit integer
type Int24 struct {
	Hi int8
	Lo uint16
}

// ToSample returns a volume as a typed value supporting the bits per sample provided
func (v Volume) ToSample(bitsPerSample int) interface{} {
	val := v.WithOverflowProtection()
	switch bitsPerSample {
	case 8:
		return int8(val * 128.0)
	case 16:
		return int16(val * 32678.0)
	case 24:
		s := int32(val * 8388608.0)
		return Int24{Hi: int8(s >> 16), Lo: uint16(s & 65535)}
	case 32:
		return int32(val * 2147483648.0)
	}
	return 0
}

// ToIntSample returns a volume as an int32 value ranged to the bits per sample provided
func (v Volume) ToIntSample(bitsPerSample int) int32 {
	val := v.WithOverflowProtection()
	switch bitsPerSample {
	case 8:
		return int32(val * 128.0)
	case 16:
		return int32(val * 32678.0)
	case 24:
		return int32(val * 8388608.0)
	case 32:
		return int32(val * 2147483648.0)
	}
	return 0
}

// ToUintSample returns a volume as a uint32 value ranged to the bits per sample provided
func (v Volume) ToUintSample(bitsPerSample int) uint32 {
	val := v.WithOverflowProtection()
	switch bitsPerSample {
	case 8:
		return uint32(val*float64(math.MaxInt8)) ^ 0x80
	case 16:
		return uint32(val*float64(math.MaxInt16)) ^ 0x8000
	case 24:
		return uint32(val*8388608.0) ^ 0x800000
	case 32:
		return uint32(val*float64(math.MaxInt32)) ^ 0x80000000
	}
	return 0
}

// Apply multiplies the volume to 1 sample, then returns the results
func (v Volume) ApplySingle(samp Volume) Volume {
	return clearDenorm(samp * v)
}

// Apply multiplies the volume to 1 sample, then returns the results
func (v Volume) ApplyMultiple(samp []Volume) []Volume {
	vols := make([]Volume, len(samp))
	for i, s := range samp {
		vols[i] = s.ApplySingle(v)
	}
	return vols
}

func (v Volume) WithOverflowProtection() float64 {
	val := float64(v)
	if math.Abs(val) <= 1.0 {
		// likely case
		return val
	} else if math.Signbit(val) {
		// overflow, negative
		return -1.0
	}
	// overflow, positive
	return 1.0
}

func (v Volume) String() string {
	var db float64
	if v > 0 {
		db = math.Log10(float64(v)) * 20.0
	} else {
		db = math.Inf(-1)
	}
	return fmt.Sprintf("%0.2f (%0.2fdB)", v, db)
}
