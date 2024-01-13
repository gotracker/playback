package types

import "github.com/gotracker/gomixing/panning"

type Panning interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
	IsInvalid() bool
	ToPosition() panning.Position
}

type PanningInformationer[TPanning Panning] interface {
	GetDefault() TPanning
	GetMax() TPanning
}

func GetPanDefault[TPanning Panning]() TPanning {
	var pd TPanning
	return any(pd).(PanningInformationer[TPanning]).GetDefault()
}

func GetPanMax[TPanning Panning]() TPanning {
	var pd TPanning
	return any(pd).(PanningInformationer[TPanning]).GetMax()
}

type PanningDeltaer[TPanning Panning] interface {
	AddDelta(d PanDelta) TPanning
}

func AddPanningDelta[TPanning Panning](v TPanning, d PanDelta) TPanning {
	return any(v).(PanningDeltaer[TPanning]).AddDelta(d)
}
