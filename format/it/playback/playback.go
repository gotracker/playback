package playback

import (
	"github.com/gotracker/gomixing/volume"
	device "github.com/gotracker/gosound"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	itFeature "github.com/gotracker/playback/format/it/feature"
	"github.com/gotracker/playback/format/it/layout"
	"github.com/gotracker/playback/format/it/pattern"
	itPeriod "github.com/gotracker/playback/format/it/period"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	playpattern "github.com/gotracker/playback/pattern"
	"github.com/gotracker/playback/player"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/output"
	"github.com/gotracker/playback/player/state"
	"github.com/gotracker/playback/song"
)

// Manager is a playback manager for IT music
type Manager struct {
	player.Tracker

	song *layout.Song

	channels  []state.ChannelState[channel.Memory, channel.Data]
	PastNotes state.PastNotesProcessor
	pattern   pattern.State

	preMixRowTxn  *playpattern.RowUpdateTransaction
	postMixRowTxn *playpattern.RowUpdateTransaction
	premix        *device.PremixData

	rowRenderState       *rowRenderState
	OnEffect             func(playback.Effect)
	longChannelOutput    bool
	enableNewNoteActions bool
}

// NewManager creates a new manager for an IT song
func NewManager(song *layout.Song) (*Manager, error) {
	m := Manager{
		Tracker: player.Tracker{
			BaseClockRate: itPeriod.ITBaseClock,
		},
		song: song,
	}

	m.PastNotes.SetMaxPerChannel(1)

	m.Tracker.Tickable = &m
	m.Tracker.Premixable = &m
	m.Tracker.Traceable = &m

	m.pattern.Reset()
	m.pattern.Orders = song.OrderList
	m.pattern.Patterns = song.Patterns

	m.SetGlobalVolume(song.Head.GlobalVolume)
	m.SetMixerVolume(song.Head.MixingVolume)

	m.SetNumChannels(len(song.ChannelSettings))
	for i, ch := range song.ChannelSettings {
		oc := m.GetOutputChannel(ch.OutputChannelNum, m.channelInit)

		cs := m.GetChannel(i)
		cs.SetSongDataInterface(song)
		cs.SetOutputChannel(oc)
		cs.SetGlobalVolume(m.GetGlobalVolume())
		cs.SetActiveVolume(ch.InitialVolume)
		cs.SetChannelVolume(ch.ChannelVolume)
		cs.SetPanEnabled(true)
		cs.SetPan(ch.InitialPanning)
		cs.SetMemory(&song.ChannelSettings[i].Memory)
		cs.SetStoredSemitone(note.UnchangedSemitone)
	}

	txn := m.pattern.StartTransaction()

	txn.Ticks.Set(song.Head.InitialSpeed)
	txn.Tempo.Set(song.Head.InitialTempo)

	if err := txn.Commit(); err != nil {
		return nil, err
	}

	return &m, nil
}

// StartPatternTransaction returns a new row update transaction for the pattern system
func (m *Manager) StartPatternTransaction() *playpattern.RowUpdateTransaction {
	return m.pattern.StartTransaction()
}

// GetNumChannels returns the number of channels
func (m *Manager) GetNumChannels() int {
	return len(m.channels)
}

func (m *Manager) semitoneSetterFactory(st note.Semitone, fn state.PeriodUpdateFunc) state.NoteOp[channel.Memory, channel.Data] {
	return doNoteCalc{
		Semitone:   st,
		UpdateFunc: fn,
	}
}

// SetNumChannels updates the song to have the specified number of channels and resets their states
func (m *Manager) SetNumChannels(num int) {
	m.channels = make([]state.ChannelState[channel.Memory, channel.Data], num)
	m.PastNotes.SetMax(channel.MaxTotalChannels - num)

	for ch := range m.channels {
		cs := &m.channels[ch]
		cs.ResetStates()
		cs.SemitoneSetterFactory = m.semitoneSetterFactory

		cs.PortaTargetPeriod.Reset()
		cs.Trigger.Reset()
		cs.RetriggerCount = 0
		_ = cs.SetData(nil)
		ocNum := m.song.GetOutputChannel(ch)
		cs.Output = m.GetOutputChannel(ocNum, m.channelInit)

		if m.enableNewNoteActions {
			cs.PastNotes = &m.PastNotes
		}
	}
}

func (m *Manager) channelInit(ch int) *output.Channel {
	return &output.Channel{
		ChannelNum:      ch,
		Filter:          nil,
		GetSampleRate:   m.GetSampleRate,
		SetGlobalVolume: m.SetGlobalVolume,
		GetOPL2Chip:     m.GetOPL2Chip,
		ChannelVolume:   volume.Volume(1),
	}
}

// SetNextOrder sets the next order index
func (m *Manager) SetNextOrder(order index.Order) error {
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
func (m *Manager) SetNextRow(row index.Row) error {
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
func (m *Manager) SetNextRowWithBacktrack(row index.Row, allowBacktrack bool) error {
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
func (m *Manager) BreakOrder() error {
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
func (m *Manager) SetTempo(tempo int) error {
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
func (m *Manager) DecreaseTempo(delta int) error {
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
func (m *Manager) IncreaseTempo(delta int) error {
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
func (m *Manager) Configure(features []feature.Feature) error {
	if err := m.Tracker.Configure(features); err != nil {
		return err
	}
	for _, feat := range features {
		switch f := feat.(type) {
		case feature.SongLoop:
			m.pattern.SongLoop = f
		case feature.PlayUntilOrderAndRow:
			m.pattern.PlayUntilOrderAndRow = f
		case itFeature.LongChannelOutput:
			m.longChannelOutput = f.Enabled
		case itFeature.NewNoteActions:
			m.enableNewNoteActions = f.Enabled
			for ch := range m.channels {
				cs := &m.channels[ch]
				if m.enableNewNoteActions {
					cs.PastNotes = &m.PastNotes
				} else {
					cs.PastNotes = nil
				}
			}
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
func (m *Manager) CanOrderLoop() bool {
	return (m.pattern.SongLoop.Count != 0)
}

// GetSongData gets the song data object
func (m *Manager) GetSongData() song.Data {
	return m.song
}

// GetChannel returns the channel interface for the specified channel number
func (m *Manager) GetChannel(ch int) *state.ChannelState[channel.Memory, channel.Data] {
	return &m.channels[ch]
}

// GetCurrentOrder returns the current order
func (m *Manager) GetCurrentOrder() index.Order {
	return m.pattern.GetCurrentOrder()
}

// GetNumOrders returns the number of orders in the song
func (m *Manager) GetNumOrders() int {
	return m.pattern.GetNumOrders()
}

// GetCurrentRow returns the current row
func (m *Manager) GetCurrentRow() index.Row {
	return m.pattern.GetCurrentRow()
}

// GetName returns the current song's name
func (m *Manager) GetName() string {
	return m.song.GetName()
}

// SetOnEffect sets the callback for an effect being generated for a channel
func (m *Manager) SetOnEffect(fn func(playback.Effect)) {
	m.OnEffect = fn
}

func (m Manager) GetOnEffect() func(playback.Effect) {
	return m.OnEffect
}

func (m *Manager) SetEnvelopePosition(v int) {
}
