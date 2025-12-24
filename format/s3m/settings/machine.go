package settings

import (
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/quirks"
)

func GetMachineSettings(modLimits bool) *settings.MachineSettings[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning] {
	if modLimits {
		return amigaMOD31Settings
	}
	return amigaS3MSettings
}

var (
	amigaMOD31Settings = quirks.GetS3MMachineSettings(quirks.ProfileST321_ModLimits, amigaVoiceFactory)
	amigaS3MSettings   = quirks.GetS3MMachineSettings(quirks.ProfileST321, amigaVoiceFactory)
)
