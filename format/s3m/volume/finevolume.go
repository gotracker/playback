package volume

import (
	"math"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/voice/types"
)

const (
	MaxFineVolume = FineVolume(0x7f)
)

type FineVolume s3mfile.Volume

var (
	_ types.VolumeMaxer[FineVolume]   = FineVolume(0)
	_ types.VolumeDeltaer[FineVolume] = FineVolume(0)
)

const finevolCoeff = volume.Volume(1) / volume.Volume(0x80)

func (v FineVolume) ToVolume() volume.Volume {
	if v != FineVolume(s3mfile.EmptyVolume) {
		return volume.Volume(min(v, MaxFineVolume)) * finevolCoeff
	}
	return volume.VolumeUseInstVol
}

func (v FineVolume) IsInvalid() bool {
	return v > 0x7f && v != FineVolume(s3mfile.EmptyVolume)
}

func (v FineVolume) IsUseInstrumentVol() bool {
	return v == FineVolume(s3mfile.EmptyVolume)
}

func (FineVolume) GetMax() FineVolume {
	return MaxFineVolume
}

func (v FineVolume) FMA(multiplier, add float32) FineVolume {
	if v == FineVolume(s3mfile.EmptyVolume) {
		return v
	}

	return min(FineVolume(max(math.FMA(float64(v), float64(multiplier), float64(add)), 0)), MaxFineVolume)
}

func (v FineVolume) AddDelta(d types.VolumeDelta) FineVolume {
	return FineVolume(min(max(int16(v)+int16(d), 0), int16(MaxFineVolume)))
}
