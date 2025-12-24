package quirks

import "github.com/gotracker/playback/player/machine/settings"

type Profile string

type Definition struct {
	Profile         Profile
	Description     string
	Quirks          settings.MachineQuirks
	MachineDefaults any
}

var registry = map[Profile]Definition{}

// Register adds or replaces a quirks definition for a profile.
func Register(def Definition) {
	registry[def.Profile] = def
}

// Resolve returns a quirks set for the requested profile. Unknown profiles fall back to an empty quirks definition.
func Resolve(profile Profile) settings.MachineQuirks {
	if def, ok := registry[profile]; ok {
		return def.Quirks
	}
	return settings.MachineQuirks{Profile: string(profile)}
}

// Get returns the full definition for a profile if registered.
func Get(profile Profile) (Definition, bool) {
	def, ok := registry[profile]
	return def, ok
}

// List returns the known quirk definitions.
func List() []Definition {
	out := make([]Definition, 0, len(registry))
	for _, def := range registry {
		out = append(out, def)
	}
	return out
}
