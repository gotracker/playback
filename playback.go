package playback

import (
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/sampler"
)

// Playback is an interface for rendering a song to output data
type Playback interface {
	Configure([]feature.Feature) error

	// runs a single tick
	//  if the onGenerate function was provided to the SetupSampler call,
	//  then the generated output will be provided through it
	Tick(s *sampler.Sampler) error

	GetNumOrders() int
	CanOrderLoop() bool
	GetName() string
}
