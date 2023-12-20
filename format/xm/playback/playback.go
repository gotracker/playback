package playback

import (
	"errors"
	"fmt"

	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback"

	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/format/xm/layout"
	"github.com/gotracker/playback/format/xm/pattern"
	xmPeriod "github.com/gotracker/playback/format/xm/period"
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

// manager is a playback manager for XM music
type manager[TPeriod period.Period] struct {
	player.Tracker

	song *layout.Song

	channels []state.ChannelState[TPeriod, channel.Memory]
	pattern  pattern.State

	preMixRowTxn  *playpattern.RowUpdateTransaction
	postMixRowTxn *playpattern.RowUpdateTransaction
	premix        *output.PremixData

	rowRenderState *rowRenderState
	OnEffect       func(playback.Effect)
}

var _ playback.Playback = (*manager[xmPeriod.Linear])(nil)
var _ playback.Playback = (*manager[xmPeriod.Amiga])(nil)
var _ playback.Channel[xmPeriod.Linear, channel.Memory] = (*state.ChannelState[xmPeriod.Linear, channel.Memory])(nil)
var _ playback.Channel[xmPeriod.Amiga, channel.Memory] = (*state.ChannelState[xmPeriod.Amiga, channel.Memory])(nil)

func (m *manager[TPeriod]) init(song *layout.Song) error {
	m.Tracker.BaseClockRate = xmPeriod.XMBaseClock
	m.song = song

	m.Tracker.Tickable = m
	m.Tracker.Premixable = m

	m.pattern.Reset()
	m.pattern.Orders = song.OrderList
	m.pattern.Patterns = song.Patterns

	m.SetGlobalVolume(song.Head.GlobalVolume)
	m.SetMixerVolume(song.Head.MixingVolume)

	m.SetNumChannels(len(song.ChannelSettings))
	for i, ch := range song.ChannelSettings {
		oc := m.GetRenderChannel(ch.OutputChannelNum, m.channelInit)

		cs := m.GetChannel(i)
		cs.SetSongDataInterface(song)
		cs.SetRenderChannel(oc)
		cs.SetGlobalVolume(m.GetGlobalVolume())
		cs.SetActiveVolume(ch.InitialVolume)
		cs.SetPanEnabled(true)
		cs.SetPan(ch.InitialPanning)
		cs.SetMemory(&song.ChannelSettings[i].Memory)
		cs.SetStoredSemitone(note.UnchangedSemitone)
	}

	txn := m.pattern.StartTransaction()

	txn.Ticks.Set(song.Head.InitialSpeed)
	txn.Tempo.Set(song.Head.InitialTempo)

	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}

// NewManager creates a new manager for an XM song
func NewManager(song *layout.Song) (playback.Playback, error) {
	if song == nil {
		return nil, errors.New("song cannot be nil")
	}

	var linearFreqSlides bool
	if len(song.ChannelSettings) != 0 {
		linearFreqSlides = song.ChannelSettings[0].Memory.Shared.LinearFreqSlides
	}

	if linearFreqSlides {
		var m manager[xmPeriod.Linear]
		if err := m.init(song); err != nil {
			return nil, fmt.Errorf("could not initialize xm linear manager: %w", err)
		}
		return &m, nil
	} else {
		var m manager[xmPeriod.Amiga]
		if err := m.init(song); err != nil {
			return nil, fmt.Errorf("could not initialize xm amiga manager: %w", err)
		}
		return &m, nil
	}
}

func (m *manager[TPeriod]) channelInit(ch int) *render.Channel {
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
func (m *manager[TPeriod]) StartPatternTransaction() *playpattern.RowUpdateTransaction {
	return m.pattern.StartTransaction()
}

// GetNumChannels returns the number of channels
func (m *manager[TPeriod]) GetNumChannels() int {
	return len(m.channels)
}

func (m *manager[TPeriod]) semitoneSetterFactory(st note.Semitone, fn state.PeriodUpdateFunc[TPeriod]) state.NoteOp[TPeriod, channel.Memory] {
	return doNoteCalc[TPeriod]{
		Semitone:   st,
		UpdateFunc: fn,
	}
}

// SetNumChannels updates the song to have the specified number of channels and resets their states
func (m *manager[TPeriod]) SetNumChannels(num int) {
	m.channels = make([]state.ChannelState[TPeriod, channel.Memory], num)

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
func (m *manager[TPeriod]) SetNextOrder(order index.Order) error {
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
func (m *manager[TPeriod]) SetNextRow(row index.Row) error {
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
func (m *manager[TPeriod]) SetNextRowWithBacktrack(row index.Row, allowBacktrack bool) error {
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
func (m *manager[TPeriod]) BreakOrder() error {
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
func (m *manager[TPeriod]) SetTempo(tempo int) error {
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
func (m *manager[TPeriod]) DecreaseTempo(delta int) error {
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
func (m *manager[TPeriod]) IncreaseTempo(delta int) error {
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
func (m *manager[TPeriod]) Configure(features []feature.Feature) error {
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
func (m *manager[TPeriod]) CanOrderLoop() bool {
	return (m.pattern.SongLoop.Count != 0)
}

// GetSongData gets the song data object
func (m *manager[TPeriod]) GetSongData() song.Data {
	return m.song
}

// GetChannel returns the channel interface for the specified channel number
func (m *manager[TPeriod]) GetChannel(ch int) *state.ChannelState[TPeriod, channel.Memory] {
	return &m.channels[ch]
}

// GetCurrentOrder returns the current order
func (m *manager[TPeriod]) GetCurrentOrder() index.Order {
	return m.pattern.GetCurrentOrder()
}

// GetNumOrders returns the number of orders in the song
func (m *manager[TPeriod]) GetNumOrders() int {
	return m.pattern.GetNumOrders()
}

// GetCurrentRow returns the current row
func (m *manager[TPeriod]) GetCurrentRow() index.Row {
	return m.pattern.GetCurrentRow()
}

// GetName returns the current song's name
func (m *manager[TPeriod]) GetName() string {
	return m.song.GetName()
}

// SetOnEffect sets the callback for an effect being generated for a channel
func (m *manager[TPeriod]) SetOnEffect(fn func(playback.Effect)) {
	m.OnEffect = fn
}

func (m *manager[TPeriod]) GetOnEffect() func(playback.Effect) {
	return m.OnEffect
}

func (m *manager[TPeriod]) SetEnvelopePosition(v int) {
}
