package playback

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/layout"
	"github.com/gotracker/playback/player"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
)

// manager is a playback manager for S3M music
type manager struct {
	player.Tracker
	song *layout.Song
}

var _ playback.Playback = (*manager)(nil)

// NewManager creates a new manager for an S3M song
func NewManager(songData song.Data) (playback.Playback, error) {
	s := songData.(*layout.Song)

	m := manager{
		song: s,
	}

	return &m, nil
}

// Configure sets specified features
func (m *manager) Configure(features []feature.Feature) error {
	us, err := m.Tracker.Configure(features)
	if err != nil {
		return err
	}
	for _, feat := range features {
		switch f := feat.(type) {
		case feature.SongLoop:
			us.SongLoop = f
		case feature.PlayUntilOrderAndRow:
			us.PlayUntilOrderAndRow = f
		case feature.SetDefaultTempo:
			us.StartTempo = f.Tempo
		case feature.SetDefaultBPM:
			us.StartBPM = f.BPM
		}
	}
	return m.SetupMachine(m.song, us)
}

// SetupSampler configures the internal sampler
func (m *manager) SetupSampler(samplesPerSecond int, channels int) error {
	if err := m.Tracker.SetupSampler(samplesPerSecond, channels); err != nil {
		return err
	}

	return nil
}
