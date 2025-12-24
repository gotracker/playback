package settings

import (
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/quirks"
)

func GetMachineSettings[TPeriod period.Period]() *settings.MachineSettings[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning] {
	var p TPeriod
	switch any(p).(type) {
	case period.Amiga:
		return any(amigaMachine).(*settings.MachineSettings[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning])
	case period.Linear:
		return any(linearMachine).(*settings.MachineSettings[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning])
	default:
		panic("unsupported machine type")
	}
}

var (
	itProfile = quirks.ProfileIT214

	amigaMachine  = quirks.GetITMachineSettingsAmiga(itProfile, amigaVoiceFactory)
	linearMachine = quirks.GetITMachineSettingsLinear(itProfile, linearVoiceFactory)
)
