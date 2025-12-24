package quirks

import (
	xmFilter "github.com/gotracker/playback/format/xm/filter"
	xmOscillator "github.com/gotracker/playback/format/xm/oscillator"
	xmPeriod "github.com/gotracker/playback/format/xm/period"
	"github.com/gotracker/playback/player/machine/settings"
)

const (
	ProfileFT210 Profile = "ft2.10"
)

func init() {
	Register(Definition{
		Profile:     ProfileFT210,
		Description: "FastTracker 2.10",
		Quirks: settings.MachineQuirks{
			Profile:                            string(ProfileFT210),
			PreviousPeriodUsesModifiedPeriod:   false,
			PortaToNoteUsesModifiedPeriod:      false,
			DoNotProcessEffectsOnMutedChannels: false,
		},
		MachineDefaults: XMMachineDefaults{
			AmigaPeriod:      xmPeriod.AmigaConverter,
			LinearPeriod:     xmPeriod.LinearConverter,
			FilterFactory:    xmFilter.Factory,
			VibratoFactory:   xmOscillator.VibratoFactory,
			TremoloFactory:   xmOscillator.TremoloFactory,
			PanbrelloFactory: xmOscillator.PanbrelloFactory,
		},
	})
}
