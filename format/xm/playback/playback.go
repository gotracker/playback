package playback

import (
	"errors"

	"github.com/gotracker/playback"

	"github.com/gotracker/playback/format/xm/layout"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
)

// manager is a playback manager for XM music
type manager[TPeriod period.Period] struct {
	player.Tracker
	song *layout.Song[TPeriod]
}

var _ playback.Playback = (*manager[period.Linear])(nil)

// NewManager creates a new manager for an XM song
func NewManager(songData song.Data) (playback.Playback, error) {
	if songData == nil {
		return nil, errors.New("song cannot be nil")
	}

	switch s := songData.(type) {
	case *layout.Song[period.Linear]:
		m := manager[period.Linear]{
			song: s,
		}
		return &m, nil

	case *layout.Song[period.Amiga]:
		m := manager[period.Amiga]{
			song: s,
		}
		return &m, nil

	default:
		return nil, errors.New("unsupported xm song data")
	}
}

// Configure sets specified features
func (m *manager[TPeriod]) Configure(features []feature.Feature) error {
	us, err := m.Tracker.Configure(features)
	if err != nil {
		return err
	}
	return m.SetupMachine(m.song, us)
}
