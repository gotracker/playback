package volume

import (
	"math"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/voice/types"
)

const (
	MaxVolume = Volume(0x40)
)

var (
	// DefaultVolume is the default volume value for most everything in S3M format
	DefaultVolume = VolumeFromS3M(Volume(s3mfile.DefaultVolume))
)

type Volume s3mfile.Volume

var (
	_ types.VolumeMaxer[Volume]   = Volume(0)
	_ types.VolumeDeltaer[Volume] = Volume(0)
)

const volCoeff = volume.Volume(1) / volume.Volume(64)

func (v Volume) ToVolume() volume.Volume {
	if v != Volume(s3mfile.EmptyVolume) {
		return volume.Volume(min(v, 0x3f)) * volCoeff
	}
	return volume.VolumeUseInstVol
}

func (v Volume) IsInvalid() bool {
	return v > 64 && v != Volume(s3mfile.EmptyVolume)
}

func (v Volume) IsUseInstrumentVol() bool {
	return v == Volume(s3mfile.EmptyVolume)
}

func (Volume) GetMax() Volume {
	return MaxVolume
}

func (v Volume) FMA(multiplier, add float32) Volume {
	if v == Volume(s3mfile.EmptyVolume) {
		return v
	}

	return Volume(min(max(math.FMA(float64(v), float64(multiplier), float64(add)), 0), float64(MaxVolume)))
}

func (v Volume) AddDelta(d types.VolumeDelta) Volume {
	return Volume(min(max(int16(v)+int16(d), 0), int16(MaxVolume)))
}

// VolumeFromS3M converts an S3M volume to a player volume
func VolumeFromS3M(vol Volume) volume.Volume {
	var v volume.Volume
	switch {
	case vol == Volume(s3mfile.EmptyVolume):
		v = volume.VolumeUseInstVol
	case vol >= 63:
		v = volume.Volume(63.0) / 64.0
	case vol < 63:
		v = volume.Volume(vol) / 64.0
	default:
		v = 0.0
	}
	return v
}

// VolumeToS3M converts a player volume to an S3M volume
func VolumeToS3M(v volume.Volume) Volume {
	switch {
	case v == volume.VolumeUseInstVol:
		return Volume(s3mfile.EmptyVolume)
	default:
		return Volume(v * 64.0)
	}
}

// VolumeFromS3M8BitSample converts an S3M 8-bit sample volume to a player volume
func VolumeFromS3M8BitSample(vol uint8) volume.Volume {
	return (volume.Volume(vol) - 128.0) / 128.0
}

// VolumeFromS3M16BitSample converts an S3M 16-bit sample volume to a player volume
func VolumeFromS3M16BitSample(vol uint16) volume.Volume {
	return (volume.Volume(vol) - 32768.0) / 32768.0
}
