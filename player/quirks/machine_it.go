package quirks

import (
	"github.com/gotracker/playback/filter"
	itFilter "github.com/gotracker/playback/format/it/filter"
	itOscillator "github.com/gotracker/playback/format/it/oscillator"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itPeriod "github.com/gotracker/playback/format/it/period"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/oscillator"
)

type ITMachineDefaults struct {
	AmigaPeriod      period.PeriodConverter[period.Amiga]
	LinearPeriod     period.PeriodConverter[period.Linear]
	FilterFactory    any
	VibratoFactory   any
	TremoloFactory   any
	PanbrelloFactory any
}

func GetITMachineSettingsAmiga(profile Profile, vf voice.VoiceFactory[period.Amiga, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) *settings.MachineSettings[period.Amiga, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning] {
	defs := getITDefaults(profile)

	return &settings.MachineSettings[period.Amiga, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]{
		PeriodConverter:     defs.AmigaPeriod,
		GetFilterFactory:    defs.FilterFactory.(func(string, frequency.Frequency, any) (filter.Filter, error)),
		GetVibratoFactory:   defs.VibratoFactory.(func() (oscillator.Oscillator, error)),
		GetTremoloFactory:   defs.TremoloFactory.(func() (oscillator.Oscillator, error)),
		GetPanbrelloFactory: defs.PanbrelloFactory.(func() (oscillator.Oscillator, error)),
		VoiceFactory:        vf,
		OPL2Enabled:         false,
		ModLimits:           false,
		Quirks:              Resolve(profile),
	}
}

func GetITMachineSettingsLinear(profile Profile, vf voice.VoiceFactory[period.Linear, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) *settings.MachineSettings[period.Linear, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning] {
	defs := getITDefaults(profile)

	return &settings.MachineSettings[period.Linear, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]{
		PeriodConverter:     defs.LinearPeriod,
		GetFilterFactory:    defs.FilterFactory.(func(string, frequency.Frequency, any) (filter.Filter, error)),
		GetVibratoFactory:   defs.VibratoFactory.(func() (oscillator.Oscillator, error)),
		GetTremoloFactory:   defs.TremoloFactory.(func() (oscillator.Oscillator, error)),
		GetPanbrelloFactory: defs.PanbrelloFactory.(func() (oscillator.Oscillator, error)),
		VoiceFactory:        vf,
		OPL2Enabled:         false,
		ModLimits:           false,
		Quirks:              Resolve(profile),
	}
}

func getITDefaults(profile Profile) ITMachineDefaults {
	if def, ok := Get(profile); ok {
		if md, ok := def.MachineDefaults.(ITMachineDefaults); ok {
			return md
		}
	}

	return ITMachineDefaults{
		AmigaPeriod:      itPeriod.AmigaConverter,
		LinearPeriod:     itPeriod.LinearConverter,
		FilterFactory:    itFilter.Factory,
		VibratoFactory:   itOscillator.VibratoFactory,
		TremoloFactory:   itOscillator.TremoloFactory,
		PanbrelloFactory: itOscillator.PanbrelloFactory,
	}
}
