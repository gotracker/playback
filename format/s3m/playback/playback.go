package playback

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/layout"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mSystem "github.com/gotracker/playback/format/s3m/system"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/player"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/system"
)

// manager is a playback manager for S3M music
type manager struct {
	player.Tracker
	song *layout.Song
}

var _ playback.Playback = (*manager)(nil)
var s3m system.ClockableSystem = s3mSystem.S3MSystem

// NewManager creates a new manager for an S3M song
func NewManager(songData song.Data) (playback.Playback, error) {
	s := songData.(*layout.Song)

	m := manager{
		Tracker: player.Tracker{
			BaseClockRate:     s3m.GetBaseClock(),
			LongChannelOutput: true,
		},
		song: s,
	}

	return &m, nil
}

func (m *manager) channelInit(ch int) render.ChannelIntf {
	return &render.Channel[s3mVolume.Volume, s3mVolume.FineVolume, s3mPanning.Panning]{
		ChannelNum:    ch,
		Filter:        nil,
		GetSampleRate: m.GetSampleRate,
		SetGlobalVolume: func(v s3mVolume.Volume) error {
			m.SetGlobalVolume(v.ToVolume())
			return nil
		},
		GetOPL2Chip:   m.GetOPL2Chip,
		ChannelVolume: s3mVolume.MaxFineVolume,
	}
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

	var oplLen int
	//oplLen += len(m.chOrder[int(s3mfile.ChannelCategoryOPL2Melody)-1])
	//oplLen += len(m.chOrder[int(s3mfile.ChannelCategoryOPL2Drums)-1])

	if oplLen > 0 {
		m.ensureOPL2()
	}
	return nil
}

func (m *manager) ensureOPL2() {
	if opl2 := m.GetOPL2Chip(); opl2 == nil {
		if s := m.GetSampler(); s != nil {
			opl2 = render.NewOPL2Chip(uint32(s.SampleRate))
			opl2.WriteReg(0x01, 0x20) // enable all waveforms
			opl2.WriteReg(0x04, 0x00) // clear timer flags
			opl2.WriteReg(0x08, 0x40) // clear CSW and set NOTE-SEL
			opl2.WriteReg(0xBD, 0x00) // set default notes
			m.SetOPL2Chip(opl2)
		}
	}
}

func (m *manager) GetNumOrders() int {
	return len(m.song.GetOrderList())
}

func (m *manager) CanOrderLoop() bool {
	return true
}

func (m *manager) GetName() string {
	return m.song.GetName()
}
