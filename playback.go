package playback

import (
	"time"

	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/player/feature"
)

// Playback is an interface for rendering a song to output data
type Playback interface {
	SetupSampler(samplesPerSecond int, channels int) error
	Configure([]feature.Feature) error

	Update(time.Duration, chan<- *output.PremixData) error
	Generate(time.Duration) (*output.PremixData, error)

	GetNumOrders() int
	CanOrderLoop() bool
	GetName() string
}
