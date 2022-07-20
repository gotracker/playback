package playback

import (
	"time"

	device "github.com/gotracker/gosound"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/pattern"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/output"
	"github.com/gotracker/playback/song"
)

// Playback is an interface for rendering a song to output data
type Playback interface {
	output.ConfigIntf

	Update(time.Duration, chan<- *device.PremixData) error
	Generate(time.Duration) (*device.PremixData, error)

	GetSongData() song.Data

	GetNumChannels() int
	GetNumOrders() int
	SetNextOrder(index.Order) error
	SetNextRow(index.Row) error
	SetNextRowWithBacktrack(index.Row, bool) error
	GetCurrentRow() index.Row
	Configure([]feature.Feature) error
	GetName() string
	CanOrderLoop() bool
	BreakOrder() error
	SetOnEffect(func(Effect))
	GetOnEffect() func(Effect)
	IgnoreUnknownEffect() bool

	StartPatternTransaction() *pattern.RowUpdateTransaction
}