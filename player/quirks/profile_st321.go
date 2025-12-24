package quirks

import (
	s3mFilter "github.com/gotracker/playback/format/s3m/filter"
	s3mOscillator "github.com/gotracker/playback/format/s3m/oscillator"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
	"github.com/gotracker/playback/player/machine/settings"
)

const (
	ProfileST321 Profile = "st3.21"
)

func init() {
	Register(Definition{
		Profile:     ProfileST321,
		Description: "Scream Tracker 3.21",
		Quirks: settings.MachineQuirks{
			Profile:                            string(ProfileST321),
			PreviousPeriodUsesModifiedPeriod:   true,
			PortaToNoteUsesModifiedPeriod:      true,
			DoNotProcessEffectsOnMutedChannels: true,
		},
		MachineDefaults: S3MMachineDefaults{
			AmigaPeriod:      s3mPeriod.S3MAmigaConverter,
			FilterFactory:    s3mFilter.Factory,
			VibratoFactory:   s3mOscillator.VibratoFactory,
			TremoloFactory:   s3mOscillator.TremoloFactory,
			PanbrelloFactory: s3mOscillator.PanbrelloFactory,
			ModLimits:        false,
		},
	})
}
