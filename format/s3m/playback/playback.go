package playback

import (
	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
	"github.com/gotracker/playback/format/s3m/layout"
	"github.com/gotracker/playback/format/s3m/pattern"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/output"
	playpattern "github.com/gotracker/playback/pattern"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/player/state"
	"github.com/gotracker/playback/song"
)

type channelState = state.ChannelState[period.Amiga, channel.Memory]

// manager is a playback manager for S3M music
type manager struct {
	player.Tracker

	song *layout.Song

	channels []channelState
	pattern  pattern.State

	preMixRowTxn  *playpattern.RowUpdateTransaction
	postMixRowTxn *playpattern.RowUpdateTransaction
	premix        *output.PremixData

	rowRenderState *rowRenderState
	OnEffect       func(playback.Effect)

	chOrder [4][]*channelState
}

var _ playback.Playback = (*manager)(nil)
var _ playback.Channel[period.Amiga, channel.Memory] = (*state.ChannelState[period.Amiga, channel.Memory])(nil)

// NewManager creates a new manager for an S3M song
func NewManager(song *layout.Song) (playback.Playback, error) {
	m := manager{
		Tracker: player.Tracker{
			BaseClockRate: s3mPeriod.S3MBaseClock,
		},
		song: song,
	}

	m.Tracker.Tickable = &m
	m.Tracker.Premixable = &m

	m.pattern.Reset()
	m.pattern.Orders = song.OrderList
	m.pattern.Patterns = song.Patterns

	m.SetGlobalVolume(song.Head.GlobalVolume)
	m.SetMixerVolume(song.Head.MixingVolume)

	m.SetNumChannels(len(song.ChannelSettings))
	lowpassEnabled := false
	for i, ch := range song.ChannelSettings {
		oc := m.GetRenderChannel(ch.OutputChannelNum, m.channelInit)

		cs := m.GetChannel(i)
		cs.PeriodConverter = s3mPeriod.AmigaConverter
		cs.SetSongDataInterface(song)
		cs.SetRenderChannel(oc)
		cs.SetGlobalVolume(m.GetGlobalVolume())
		cs.SetActiveVolume(ch.InitialVolume)
		if song.Head.Stereo {
			cs.SetPanEnabled(true)
			cs.SetPan(ch.InitialPanning)
		} else {
			cs.SetPanEnabled(true)
			cs.SetPan(panning.CenterAhead)
			cs.SetPanEnabled(false)
		}
		cs.SetStoredSemitone(note.UnchangedSemitone)
		mem := &song.ChannelSettings[i].Memory
		cs.SetMemory(mem)
		if mem.Shared.LowPassFilterEnable {
			lowpassEnabled = true
		}

		// weirdly, S3M processes channels in channel category order
		// so we have to make a list with the order we're expecting
		switch s3mfile.ChannelCategory(ch.Category) {
		case s3mfile.ChannelCategoryUnknown:
			// do nothing
		default:
			cIdx := int(ch.Category) - 1
			m.chOrder[cIdx] = append(m.chOrder[cIdx], cs)
		}
	}

	if lowpassEnabled {
		m.SetFilterEnable(true)
	}

	txn := m.pattern.StartTransaction()

	txn.Ticks.Set(song.Head.InitialSpeed)
	txn.Tempo.Set(song.Head.InitialTempo)

	if err := txn.Commit(); err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *manager) channelInit(ch int) *render.Channel {
	return &render.Channel{
		ChannelNum:      ch,
		Filter:          nil,
		GetSampleRate:   m.GetSampleRate,
		SetGlobalVolume: m.SetGlobalVolume,
		GetOPL2Chip:     m.GetOPL2Chip,
		ChannelVolume:   volume.Volume(1),
	}
}

// StartPatternTransaction returns a new row update transaction for the pattern system
func (m *manager) StartPatternTransaction() *playpattern.RowUpdateTransaction {
	return m.pattern.StartTransaction()
}

// GetNumChannels returns the number of channels
func (m *manager) GetNumChannels() int {
	return len(m.channels)
}

func (m *manager) semitoneSetterFactory(st note.Semitone, fn state.PeriodUpdateFunc[period.Amiga]) state.NoteOp[period.Amiga, channel.Memory] {
	return doNoteCalc{
		Semitone:   st,
		UpdateFunc: fn,
	}
}

// SetNumChannels updates the song to have the specified number of channels and resets their states
func (m *manager) SetNumChannels(num int) {
	m.channels = make([]channelState, num)

	for ch := range m.channels {
		cs := &m.channels[ch]
		cs.ResetStates()
		cs.SemitoneSetterFactory = m.semitoneSetterFactory

		cs.PortaTargetPeriod.Reset()
		cs.Trigger.Reset()
		cs.RetriggerCount = 0
		_ = cs.SetData(nil)
		ocNum := m.song.GetRenderChannel(ch)
		cs.RenderChannel = m.GetRenderChannel(ocNum, m.channelInit)
	}
}

// SetNextOrder sets the next order index
func (m *manager) SetNextOrder(order index.Order) error {
	if m.postMixRowTxn != nil {
		m.postMixRowTxn.SetNextOrder(order)
	} else {
		rowTxn := m.pattern.StartTransaction()
		defer rowTxn.Cancel()

		rowTxn.SetNextOrder(order)
		if err := rowTxn.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// SetNextRow sets the next row index
func (m *manager) SetNextRow(row index.Row) error {
	if m.postMixRowTxn != nil {
		m.postMixRowTxn.SetNextRow(row)
	} else {
		rowTxn := m.pattern.StartTransaction()
		defer rowTxn.Cancel()

		rowTxn.SetNextRow(row)
		if err := rowTxn.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// SetNextRowWithBacktrack will set the next row index and backtracing allowance
func (m *manager) SetNextRowWithBacktrack(row index.Row, allowBacktrack bool) error {
	if m.postMixRowTxn != nil {
		m.postMixRowTxn.SetNextRowWithBacktrack(row, allowBacktrack)
	} else {
		rowTxn := m.pattern.StartTransaction()
		defer rowTxn.Cancel()

		rowTxn.SetNextRowWithBacktrack(row, allowBacktrack)
		if err := rowTxn.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// BreakOrder breaks to the next pattern in the order
func (m *manager) BreakOrder() error {
	if m.postMixRowTxn != nil {
		m.postMixRowTxn.BreakOrder = true
	} else {
		rowTxn := m.pattern.StartTransaction()
		defer rowTxn.Cancel()

		rowTxn.BreakOrder = true
		if err := rowTxn.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// SetTempo sets the desired tempo for the song
func (m *manager) SetTempo(tempo int) error {
	if m.preMixRowTxn != nil {
		m.preMixRowTxn.Tempo.Set(tempo)
	} else {
		rowTxn := m.pattern.StartTransaction()
		defer rowTxn.Cancel()

		rowTxn.Tempo.Set(tempo)
		if err := rowTxn.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// DecreaseTempo reduces the tempo by the `delta` value
func (m *manager) DecreaseTempo(delta int) error {
	if m.preMixRowTxn != nil {
		m.preMixRowTxn.AccTempoDelta(-delta)
	} else {
		rowTxn := m.pattern.StartTransaction()
		defer rowTxn.Cancel()

		rowTxn.AccTempoDelta(-delta)
		if err := rowTxn.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// IncreaseTempo increases the tempo by the `delta` value
func (m *manager) IncreaseTempo(delta int) error {
	if m.preMixRowTxn != nil {
		m.preMixRowTxn.AccTempoDelta(delta)
	} else {
		rowTxn := m.pattern.StartTransaction()
		defer rowTxn.Cancel()

		rowTxn.AccTempoDelta(delta)
		if err := rowTxn.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// Configure sets specified features
func (m *manager) Configure(features []feature.Feature) error {
	if err := m.Tracker.Configure(features); err != nil {
		return err
	}
	for _, feat := range features {
		switch f := feat.(type) {
		case feature.SongLoop:
			m.pattern.SongLoop = f
		case feature.PlayUntilOrderAndRow:
			m.pattern.PlayUntilOrderAndRow = f
		case feature.SetDefaultTempo:
			txn := m.pattern.StartTransaction()
			txn.Ticks.Set(f.Tempo)
			if err := txn.Commit(); err != nil {
				return err
			}
		case feature.SetDefaultBPM:
			txn := m.pattern.StartTransaction()
			txn.Tempo.Set(f.BPM)
			if err := txn.Commit(); err != nil {
				return err
			}
		}
	}
	return nil
}

// CanOrderLoop returns true if the song is allowed to order loop
func (m *manager) CanOrderLoop() bool {
	return (m.pattern.SongLoop.Count != 0)
}

// GetSongData gets the song data object
func (m *manager) GetSongData() song.Data {
	return m.song
}

// GetChannel returns the channel interface for the specified channel number
func (m *manager) GetChannel(ch int) *channelState {
	return &m.channels[ch]
}

// GetCurrentOrder returns the current order
func (m *manager) GetCurrentOrder() index.Order {
	return m.pattern.GetCurrentOrder()
}

// GetNumOrders returns the number of orders in the song
func (m *manager) GetNumOrders() int {
	return m.pattern.GetNumOrders()
}

// GetCurrentRow returns the current row
func (m *manager) GetCurrentRow() index.Row {
	return m.pattern.GetCurrentRow()
}

// GetName returns the current song's name
func (m *manager) GetName() string {
	return m.song.GetName()
}

// SetOnEffect sets the callback for an effect being generated for a channel
func (m *manager) SetOnEffect(fn func(playback.Effect)) {
	m.OnEffect = fn
}

func (m *manager) GetOnEffect() func(playback.Effect) {
	return m.OnEffect
}

func (m *manager) SetEnvelopePosition(v int) {
}

// SetupSampler configures the internal sampler
func (m *manager) SetupSampler(samplesPerSecond int, channels int) error {
	if err := m.Tracker.SetupSampler(samplesPerSecond, channels); err != nil {
		return err
	}

	oplLen := len(m.chOrder[int(s3mfile.ChannelCategoryOPL2Melody)-1])
	oplLen += len(m.chOrder[int(s3mfile.ChannelCategoryOPL2Drums)-1])

	if oplLen > 0 {
		m.ensureOPL2()
	}
	return nil
}
