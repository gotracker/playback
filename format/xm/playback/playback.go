package playback

import (
	"errors"
	"fmt"

	"github.com/gotracker/playback"

	"github.com/gotracker/playback/format/xm/layout"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmSystem "github.com/gotracker/playback/format/xm/system"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/song"
)

// manager is a playback manager for XM music
type manager[TPeriod period.Period] struct {
	player.Tracker
	song *layout.Song[TPeriod]
}

var _ playback.Playback = (*manager[period.Linear])(nil)
var _ playback.Playback = (*manager[period.Amiga])(nil)

func (m *manager[TPeriod]) init(s *layout.Song[TPeriod]) error {
	m.Tracker.BaseClockRate = xmSystem.XMBaseClock
	m.Tracker.LongChannelOutput = true
	m.song = s

	return nil
}

// NewManager creates a new manager for an XM song
func NewManager(songData song.Data) (playback.Playback, error) {
	if songData == nil {
		return nil, errors.New("song cannot be nil")
	}

	switch s := songData.(type) {
	case *layout.Song[period.Linear]:
		var m manager[period.Linear]
		if err := m.init(s); err != nil {
			return nil, fmt.Errorf("could not initialize xm manager: %w", err)
		}
		return &m, nil

	case *layout.Song[period.Amiga]:
		var m manager[period.Amiga]
		if err := m.init(s); err != nil {
			return nil, fmt.Errorf("could not initialize xm manager: %w", err)
		}
		return &m, nil

	default:
		return nil, errors.New("unsupported xm song data")
	}
}

func (m *manager[TPeriod]) channelInit(ch int) render.ChannelIntf {
	return &render.Channel[xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]{
		ChannelNum:    ch,
		Filter:        nil,
		GetSampleRate: m.GetSampleRate,
		SetGlobalVolume: func(xv xmVolume.XmVolume) error {
			m.SetGlobalVolume(xv.ToVolume())
			return nil
		},
		GetOPL2Chip:   m.GetOPL2Chip,
		ChannelVolume: xmVolume.DefaultXmVolume,
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

func (m *manager[TPeriod]) GetNumOrders() int {
	return len(m.song.GetOrderList())
}

func (m *manager[TPeriod]) CanOrderLoop() bool {
	return true
}

func (m *manager[TPeriod]) GetName() string {
	return m.song.GetName()
}
