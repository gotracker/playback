package settings

import (
	xmFilter "github.com/gotracker/playback/format/xm/filter"
	xmOscillator "github.com/gotracker/playback/format/xm/oscillator"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmPeriod "github.com/gotracker/playback/format/xm/period"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
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
		GetFilterFactory:    xmFilter.Factory,
		GetVibratoFactory:   xmOscillator.VibratoFactory,
		GetTremoloFactory:   xmOscillator.TremoloFactory,
		GetPanbrelloFactory: xmOscillator.PanbrelloFactory,
		VoiceFactory:        amigaVoiceFactory,
		OPL2Enabled:         false,
	}

	linearMachine = settings.MachineSettings[period.Linear, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]{
		PeriodConverter:     xmPeriod.LinearConverter,
		GetFilterFactory:    xmFilter.Factory,
		GetVibratoFactory:   xmOscillator.VibratoFactory,
		GetTremoloFactory:   xmOscillator.TremoloFactory,
		GetPanbrelloFactory: xmOscillator.PanbrelloFactory,
		VoiceFactory:        linearVoiceFactory,
		OPL2Enabled:         false,
	}
)
