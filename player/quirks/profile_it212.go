package quirks

import (
	itFilter "github.com/gotracker/playback/format/it/filter"
	itOscillator "github.com/gotracker/playback/format/it/oscillator"
	itPeriod "github.com/gotracker/playback/format/it/period"
	"github.com/gotracker/playback/player/machine/settings"
)

const (
	ProfileIT212 Profile = "it212"
)

func init() {
	Register(Definition{
		Profile:     ProfileIT212,
		Description: "Impulse Tracker 2.12 (older behavior)",
		Quirks: settings.MachineQuirks{
			Profile:                            string(ProfileIT212),
			PreviousPeriodUsesModifiedPeriod:   false,
			PortaToNoteUsesModifiedPeriod:      false,
			DoNotProcessEffectsOnMutedChannels: false,
		},
		MachineDefaults: ITMachineDefaults{
			AmigaPeriod:      itPeriod.AmigaConverter,
			LinearPeriod:     itPeriod.LinearConverter,
			FilterFactory:    itFilter.Factory,
			VibratoFactory:   itOscillator.VibratoFactory,
			TremoloFactory:   itOscillator.TremoloFactory,
			PanbrelloFactory: itOscillator.PanbrelloFactory,
		},
	})
}
