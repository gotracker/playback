package quirks

import (
	itFilter "github.com/gotracker/playback/format/it/filter"
	itOscillator "github.com/gotracker/playback/format/it/oscillator"
	itPeriod "github.com/gotracker/playback/format/it/period"
	"github.com/gotracker/playback/player/machine/settings"
)

const (
	ProfileIT214 Profile = "it214"
)

func init() {
	Register(Definition{
		Profile:     ProfileIT214,
		Description: "Impulse Tracker 2.14 (classic behavior)",
		Quirks: settings.MachineQuirks{
			Profile:                            string(ProfileIT214),
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
