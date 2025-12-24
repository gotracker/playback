package quirks

import "github.com/gotracker/playback/player/machine/settings"

const (
	ProfileOpenMPTCurrent Profile = "openmpt-current"
)

func init() {
	Register(Definition{
		Profile:     ProfileOpenMPTCurrent,
		Description: "OpenMPT (modern defaults)",
		Quirks: settings.MachineQuirks{
			Profile:                            string(ProfileOpenMPTCurrent),
			PreviousPeriodUsesModifiedPeriod:   false,
			PortaToNoteUsesModifiedPeriod:      false,
			DoNotProcessEffectsOnMutedChannels: false,
		},
	})
}
