package settings

import (
	"fmt"

	"github.com/gotracker/playback/filter"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itPeriod "github.com/gotracker/playback/format/it/period"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/oscillator"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/player/machine/settings"
	voiceOscillator "github.com/gotracker/playback/voice/oscillator"
)

func GetMachineSettings[TPeriod period.Period]() *settings.MachineSettings[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning] {
	var p TPeriod
	switch any(p).(type) {
	case period.Amiga:
		return any(&amigaMachine).(*settings.MachineSettings[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning])
	case period.Linear:
		return any(&linearMachine).(*settings.MachineSettings[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning])
	default:
		panic("unsupported machine type")
	}
}

var (
	amigaMachine = settings.MachineSettings[period.Amiga, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]{
		PeriodConverter:     itPeriod.AmigaConverter,
		GetFilterFactory:    filterFactory,
		GetVibratoFactory:   vibratoFactory,
		GetTremoloFactory:   tremoloFactory,
		GetPanbrelloFactory: panbrelloFactory,
		VoiceFactory:        amigaVoiceFactory,
		OPL2Enabled:         false,
	}

	linearMachine = settings.MachineSettings[period.Linear, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]{
		PeriodConverter:     itPeriod.LinearConverter,
		GetFilterFactory:    filterFactory,
		GetVibratoFactory:   vibratoFactory,
		GetTremoloFactory:   tremoloFactory,
		GetPanbrelloFactory: panbrelloFactory,
		VoiceFactory:        linearVoiceFactory,
		OPL2Enabled:         false,
	}
)

func filterFactory(name string) (settings.FilterFactoryFunc, error) {
	switch name {
	case "amigalpf":
		return func(instrument, playback period.Frequency) (filter.Filter, error) {
			lpf := filter.NewAmigaLPF(instrument, playback)
			return lpf, nil
		}, nil

	default:
		return nil, fmt.Errorf("unsupported filter: %q", name)
	}
}

func vibratoFactory() (voiceOscillator.Oscillator, error) {
	return oscillator.NewImpulseTrackerOscillator(4), nil
}

func tremoloFactory() (voiceOscillator.Oscillator, error) {
	return oscillator.NewImpulseTrackerOscillator(4), nil
}

func panbrelloFactory() (voiceOscillator.Oscillator, error) {
	return oscillator.NewImpulseTrackerOscillator(1), nil
}

func init() {
	machine.RegisterMachine(GetMachineSettings[period.Amiga]())
	machine.RegisterMachine(GetMachineSettings[period.Linear]())
}
