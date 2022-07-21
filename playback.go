package playback

import (
	"time"

	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/pattern"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice/render"
)

// Playback is an interface for rendering a song to output data
type Playback interface {
	SetupSampler(int, int, int) error
	GetSampleRate() period.Frequency
	GetOPL2Chip() render.OPL2Chip
	GetGlobalVolume() volume.Volume
	SetGlobalVolume(volume.Volume)

	Update(time.Duration, chan<- *output.PremixData) error
	Generate(time.Duration) (*output.PremixData, error)

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
