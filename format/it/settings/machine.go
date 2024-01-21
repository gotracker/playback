package settings

import (
	itFilter "github.com/gotracker/playback/format/it/filter"
	itOscillator "github.com/gotracker/playback/format/it/oscillator"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itPeriod "github.com/gotracker/playback/format/it/period"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
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
		GetFilterFactory:    itFilter.Factory,
		GetVibratoFactory:   itOscillator.VibratoFactory,
		GetTremoloFactory:   itOscillator.TremoloFactory,
		GetPanbrelloFactory: itOscillator.PanbrelloFactory,
		VoiceFactory:        amigaVoiceFactory,
		OPL2Enabled:         false,
	}

	linearMachine = settings.MachineSettings[period.Linear, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]{
		PeriodConverter:     itPeriod.LinearConverter,
		GetFilterFactory:    itFilter.Factory,
		GetVibratoFactory:   itOscillator.VibratoFactory,
		GetTremoloFactory:   itOscillator.TremoloFactory,
		GetPanbrelloFactory: itOscillator.PanbrelloFactory,
		VoiceFactory:        linearVoiceFactory,
		OPL2Enabled:         false,
	}
)
