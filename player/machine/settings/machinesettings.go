package settings

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/oscillator"
)

type (
	Period  = voice.Period
	Volume  = voice.Volume
	Panning = voice.Panning
)

type FilterFactoryFunc func(instrument, playback period.Frequency) (filter.Filter, error)

type MachineSettings[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] struct {
	PeriodConverter     song.PeriodCalculator[TPeriod]
	GetFilterFactory    func(name string) (FilterFactoryFunc, error)
	GetVibratoFactory   func() (oscillator.Oscillator, error)
	GetTremoloFactory   func() (oscillator.Oscillator, error)
	GetPanbrelloFactory func() (oscillator.Oscillator, error)
	VoiceFactory        voice.VoiceFactory[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]
	OPL2Enabled         bool
}
