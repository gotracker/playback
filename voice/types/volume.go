package types

import "github.com/gotracker/gomixing/volume"

type Volume interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
	IsInvalid() bool
	IsUseInstrumentVol() bool
	ToVolume() volume.Volume
}

type VolumeMaxer[TVolume Volume] interface {
	GetMax() TVolume
}

func GetMaxVolume[TVolume Volume]() TVolume {
	var vm TVolume
	return any(vm).(VolumeMaxer[TVolume]).GetMax()
}

type VolumeDeltaer[TVolume Volume] interface {
	AddDelta(d VolumeDelta) TVolume
}

func AddVolumeDelta[TVolume Volume](v TVolume, d VolumeDelta) TVolume {
	return any(v).(VolumeDeltaer[TVolume]).AddDelta(d)
}
