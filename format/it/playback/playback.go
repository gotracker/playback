package playback

import (
	"errors"

	"github.com/gotracker/playback"
	itFeature "github.com/gotracker/playback/format/it/feature"
	"github.com/gotracker/playback/format/it/layout"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
)

// manager is a playback manager for IT music
type manager[TPeriod period.Period] struct {
	player.Tracker
	song *layout.Song[TPeriod]
}

var _ playback.Playback = (*manager[period.Linear])(nil)

// NewManager creates a new manager for an IT song
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
		return nil, errors.New("unsupported it song data")
	}
}

// Configure sets specified features
func (m *manager[TPeriod]) Configure(features []feature.Feature) error {
	us, err := m.Tracker.Configure(features)
	if err != nil {
		return err
	}
	for _, feat := range features {
		switch f := feat.(type) {
		case feature.SongLoop:
			us.SongLoop = f
		case feature.StartOrderAndRow:
			us.StartOrderAndRow = f
		case feature.PlayUntilOrderAndRow:
			us.PlayUntilOrderAndRow = f
		case itFeature.LongChannelOutput:
			us.LongChannelOutput = f.Enabled
		case itFeature.NewNoteActions:
			us.EnableNewNoteActions = f.Enabled
		case feature.SetDefaultTempo:
			us.StartTempo = f.Tempo
		case feature.SetDefaultBPM:
			us.StartBPM = f.BPM
		}
	}
	return m.SetupMachine(m.song, us)
}
