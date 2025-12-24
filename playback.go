package playback

import (
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/sampler"
)

// Playback is an interface for rendering a song to output data
type Playback interface {
	Configure([]feature.Feature) error

	// Tick is a convenience for Advance+Render; sampler must be non-nil.
	Tick(s *sampler.Sampler) error

	// Advance progresses sequencing without rendering audio.
	Advance() error

	// Render produces audio for the current state using the provided sampler.
	Render(s *sampler.Sampler) error

	GetNumOrders() int
	CanOrderLoop() bool
	GetName() string
}
