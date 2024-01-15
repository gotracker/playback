package settings

import (
	"fmt"

	"github.com/gotracker/playback/filter"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmPeriod "github.com/gotracker/playback/format/xm/period"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/frequency"
	oscillatorImpl "github.com/gotracker/playback/oscillator"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/voice/oscillator"
)

func GetMachineSettings[TPeriod period.Period]() *settings.MachineSettings[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning] {
	var p TPeriod
	switch any(p).(type) {
	case period.Amiga:
		return any(&amigaMachine).(*settings.MachineSettings[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning])
	case period.Linear:
		return any(&linearMachine).(*settings.MachineSettings[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning])
	default:
		panic("unsupported machine type")
	}
}

var (
	amigaMachine = settings.MachineSettings[period.Amiga, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]{
		PeriodConverter:     xmPeriod.AmigaConverter,
		GetFilterFactory:    filterFactory,
		GetVibratoFactory:   vibratoFactory,
		GetTremoloFactory:   tremoloFactory,
		GetPanbrelloFactory: panbrelloFactory,
		VoiceFactory:        amigaVoiceFactory,
		OPL2Enabled:         false,
	}

	linearMachine = settings.MachineSettings[period.Linear, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]{
		PeriodConverter:     xmPeriod.LinearConverter,
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
		return func(instrument frequency.Frequency) (filter.Filter, error) {
			lpf := filter.NewAmigaLPF(instrument)
			return lpf, nil
		}, nil

	default:
		return nil, fmt.Errorf("unsupported filter: %q", name)
	}
}

func vibratoFactory() (oscillator.Oscillator, error) {
	return oscillatorImpl.NewImpulseTrackerOscillator(4), nil
}

func tremoloFactory() (oscillator.Oscillator, error) {
	return oscillatorImpl.NewImpulseTrackerOscillator(4), nil
}

func panbrelloFactory() (oscillator.Oscillator, error) {
	return oscillatorImpl.NewImpulseTrackerOscillator(1), nil
}
