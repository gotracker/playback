package machine

import (
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/quirks"
)

func resolveQuirks(base settings.MachineQuirks, us settings.UserSettings) settings.MachineQuirks {
	// start from the machine's preferred profile so defaults always come from the profile definitions
	if base.Profile != "" {
		base = quirks.Resolve(quirks.Profile(base.Profile))
	}

	q := base
	customized := false

	if prof, ok := us.Quirks.Profile.Get(); ok {
		q = quirks.Resolve(quirks.Profile(prof))
		customized = true
	}

	if value, ok := us.Quirks.PreviousPeriodUsesModifiedPeriodOverride.Get(); ok {
		q.PreviousPeriodUsesModifiedPeriod = value
		customized = true
	}
	if value, ok := us.Quirks.PortaToNoteUsesModifiedPeriodOverride.Get(); ok {
		q.PortaToNoteUsesModifiedPeriod = value
		customized = true
	}
	if value, ok := us.Quirks.DoNotProcessEffectsOnMutedChannelsOverride.Get(); ok {
		q.DoNotProcessEffectsOnMutedChannels = value
		customized = true
	}

	// keep profile label meaningful after overrides
	if prof, ok := us.Quirks.Profile.Get(); ok {
		q.Profile = prof
	} else if customized {
		if base.Profile != "" {
			q.Profile = base.Profile + "+custom"
		} else {
			q.Profile = "custom"
		}
	}

	return q
}
