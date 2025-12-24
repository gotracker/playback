package feature

import "github.com/heucuva/optional"

// QuirksMode carries overrides for tracker quirks handling.
// Profile is optional and maps to a tracker profile name; LinearSlides toggles
// whether linear frequency slides should be used regardless of the file flag.
type QuirksMode struct {
	Profile      optional.Value[string]
	LinearSlides optional.Value[bool]
}
