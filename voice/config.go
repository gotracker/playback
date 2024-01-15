package voice

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/voice/autovibrato"
	"github.com/gotracker/playback/voice/envelope"
	"github.com/gotracker/playback/voice/fadeout"
	"github.com/gotracker/playback/voice/types"
	"github.com/gotracker/playback/voice/vol0optimization"
)

type (
	Period  = types.Period
	Volume  = types.Volume
	Panning = types.Panning
)

type InstrumentConfig[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] struct {
	SampleRate           frequency.Frequency
	AutoVibrato          autovibrato.AutoVibratoSettings
	Data                 instrument.Data
	VoiceFilter          filter.Filter
	FadeOut              fadeout.Settings
	PitchPan             instrument.PitchPan
	VolEnv               envelope.Envelope[TVolume]
	VolEnvFinishFadesOut bool
	PanEnv               envelope.Envelope[TPanning]
	PitchFiltMode        bool                                    // true = filter, false = pitch
	PitchFiltEnv         envelope.Envelope[types.PitchFiltValue] // this is either pitch or filter
}

type VoiceConfig[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] struct {
	InitialVolume    TVolume
	InitialMixing    TMixingVolume
	PanEnabled       bool
	InitialPan       TPanning
	Vol0Optimization vol0optimization.Vol0OptimizationSettings
}
