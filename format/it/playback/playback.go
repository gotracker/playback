package playback

import (
	"errors"
	"fmt"

	"github.com/gotracker/playback"
	itFeature "github.com/gotracker/playback/format/it/feature"
	"github.com/gotracker/playback/format/it/layout"
	itSystem "github.com/gotracker/playback/format/it/system"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/system"
)

// manager is a playback manager for IT music
type manager[TPeriod period.Period] struct {
	player.Tracker
	song *layout.Song[TPeriod]
}

var _ playback.Playback = (*manager[period.Linear])(nil)

var it system.ClockableSystem = itSystem.ITSystem

func (m *manager[TPeriod]) init(s *layout.Song[TPeriod]) error {
	m.Tracker.BaseClockRate = it.GetBaseClock()
	m.song = s

	return nil
}

// NewManager creates a new manager for an IT song
func NewManager(songData song.Data) (playback.Playback, error) {
	if songData == nil {
		return nil, errors.New("song cannot be nil")
	}

	switch s := songData.(type) {
	case *layout.Song[period.Linear]:
		var m manager[period.Linear]
		if err := m.init(s); err != nil {
			return nil, fmt.Errorf("could not initialize it manager: %w", err)
		}
		return &m, nil
	case *layout.Song[period.Amiga]:
		var m manager[period.Amiga]
		if err := m.init(s); err != nil {
			return nil, fmt.Errorf("could not initialize it manager: %w", err)
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
			m.Tracker.LongChannelOutput = f.Enabled
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

func (m *manager[TPeriod]) GetNumOrders() int {
	return len(m.song.GetOrderList())
}

func (m *manager[TPeriod]) CanOrderLoop() bool {
	return true
}

func (m *manager[TPeriod]) GetName() string {
	return m.song.GetName()
}
