package settings

import (
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/quirks"
)

func GetMachineSettings[TPeriod period.Period]() *settings.MachineSettings[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning] {
	var p TPeriod
	switch any(p).(type) {
	case period.Amiga:
		return any(amigaMachine).(*settings.MachineSettings[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning])
	case period.Linear:
		return any(linearMachine).(*settings.MachineSettings[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning])
	default:
		panic("unsupported machine type")
	}
}

var (
	xmProfile = quirks.ProfileFT210

	amigaMachine  = quirks.GetXMMachineSettingsAmiga(xmProfile, amigaVoiceFactory)
	linearMachine = quirks.GetXMMachineSettingsLinear(xmProfile, linearVoiceFactory)
)
