package common

import "github.com/gotracker/playback/player/feature"

// ResolveLinearSlides returns the desired linear slides setting, allowing
// user-provided quirks configuration to override the file's default flag.
func ResolveLinearSlides(defaultLinear bool, features []feature.Feature) bool {
	for _, feat := range features {
		if qm, ok := feat.(feature.QuirksMode); ok {
			if linear, ok := qm.LinearSlides.Get(); ok {
				return linear
			}
		}
	}
	return defaultLinear
}
