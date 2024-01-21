package voice

import (
	"github.com/gotracker/opl2"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/types"
	"github.com/gotracker/playback/voice/vol0optimization"
)

type (
	Period  = types.Period
	Volume  = types.Volume
	Panning = types.Panning
)

type VoiceConfig[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] struct {
	PC               period.PeriodConverter[TPeriod]
	OPLChip          *opl2.Chip
	OPLChannel       index.OPLChannel
	InitialVolume    TVolume
	InitialMixing    TMixingVolume
	PanEnabled       bool
	InitialPan       TPanning
	Vol0Optimization vol0optimization.Vol0OptimizationSettings
}
