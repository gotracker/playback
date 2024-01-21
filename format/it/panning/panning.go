package panning

import (
	"math"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/playback/voice/types"
)

var (
	DefaultPanningLeft = Panning(0x30)
	// DefaultPanningLeftPosition is the default panning value for left channels
	DefaultPanningLeftPosition = FromItPanning(itfile.PanValue(DefaultPanningLeft))

	DefaultPanning = Panning(0x80)
	// DefaultPanningPosition is the default panning value for unconfigured channels
	DefaultPanningPosition = FromItPanning(itfile.PanValue(DefaultPanning))

	DefaultPanningRight = Panning(0xC0)
	// DefaultPanningRightPosition is the default panning value for right channels
	DefaultPanningRightPosition = FromItPanning(itfile.PanValue(DefaultPanningRight))

	MaxPanning = Panning(0xFF)
)

type Panning uint8

var (
	_ types.PanningInformationer[Panning] = Panning(0)
	_ types.PanningDeltaer[Panning]       = Panning(0)
)

func (p Panning) IsInvalid() bool {
	return false
}

func (p Panning) ToPosition() panning.Position {
	return panning.MakeStereoPosition(float32(p), 0, 0xFF)
}

func (Panning) GetDefault() Panning {
	return DefaultPanning
}

func (Panning) GetMax() Panning {
	return MaxPanning
}

func (p Panning) FMA(multiplier, add float32) Panning {
	return Panning(min(max(math.FMA(float64(p), float64(multiplier), float64(add)), 0), 0xFF))
}

func (p Panning) AddDelta(d types.PanDelta) Panning {
	return Panning(min(max(int16(p)+int16(d), 0), int16(MaxPanning)))
}

// FromItPanning returns a radian panning position from an it panning value
func FromItPanning(pos itfile.PanValue) panning.Position {
	if pos.IsDisabled() {
		return panning.CenterAhead
	}
	return panning.MakeStereoPosition(pos.Value(), 0, 1)
}

// ToItPanning returns the it panning value for a radian panning position
func ToItPanning(pan panning.Position) itfile.PanValue {
	p := panning.FromStereoPosition(pan, 0, 1)
	return itfile.PanValue(p * 64)
}
