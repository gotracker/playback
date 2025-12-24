package quirks

import (
	"github.com/gotracker/playback/filter"
	xmFilter "github.com/gotracker/playback/format/xm/filter"
	xmOscillator "github.com/gotracker/playback/format/xm/oscillator"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmPeriod "github.com/gotracker/playback/format/xm/period"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/oscillator"
)

type XMMachineDefaults struct {
	AmigaPeriod      period.PeriodConverter[period.Amiga]
	LinearPeriod     period.PeriodConverter[period.Linear]
	FilterFactory    any
	VibratoFactory   any
	TremoloFactory   any
	PanbrelloFactory any
}

func GetXMMachineSettingsAmiga(profile Profile, vf voice.VoiceFactory[period.Amiga, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) *settings.MachineSettings[period.Amiga, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning] {
	defs := getXMDefaults(profile)

	return &settings.MachineSettings[period.Amiga, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]{
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

func GetXMMachineSettingsLinear(profile Profile, vf voice.VoiceFactory[period.Linear, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) *settings.MachineSettings[period.Linear, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning] {
	defs := getXMDefaults(profile)

	return &settings.MachineSettings[period.Linear, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]{
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

func getXMDefaults(profile Profile) XMMachineDefaults {
	if def, ok := Get(profile); ok {
		if md, ok := def.MachineDefaults.(XMMachineDefaults); ok {
			return md
		}
	}

	return XMMachineDefaults{
		AmigaPeriod:      xmPeriod.AmigaConverter,
		LinearPeriod:     xmPeriod.LinearConverter,
		FilterFactory:    xmFilter.Factory,
		VibratoFactory:   xmOscillator.VibratoFactory,
		TremoloFactory:   xmOscillator.TremoloFactory,
		PanbrelloFactory: xmOscillator.PanbrelloFactory,
	}
}
