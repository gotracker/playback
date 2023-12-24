package playback

import (
	"errors"
	"time"

	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/player/state"
	"github.com/gotracker/playback/song"
)

const (
	tickBaseDuration = time.Duration(2500) * time.Millisecond
)

func (m *manager[TPeriod]) processPatternRow() error {
	patIdx, err := m.pattern.GetCurrentPatternIdx()
	if err != nil {
		return err
	}

	if m.pattern.NeedResetPatternLoops() {
		for _, cs := range m.channels {
			mem := cs.GetMemory()
			pl := mem.GetPatternLoop()
			pl.Count = 0
			pl.Enabled = false
		}
	}

	pat := m.song.GetPattern(patIdx)
	if pat == nil {
		return song.ErrStopSong
	}

	withinPatternLoop := false
	for _, cs := range m.channels {
		mem := cs.GetMemory()
		pl := mem.GetPatternLoop()
		if pl.Enabled {
			withinPatternLoop = true
			break
		}
	}

	if !withinPatternLoop {
		if err := m.pattern.Observe(); err != nil {
			return err
		}
	}

	myCurrentRow := m.pattern.GetCurrentRow()

	row := pat.GetRow(myCurrentRow)

	preMixRowTxn := m.pattern.StartTransaction()
	defer func() {
		preMixRowTxn.Cancel()
		m.preMixRowTxn = nil
	}()
	m.preMixRowTxn = preMixRowTxn

	s := m.GetSampler()
	if s == nil {
		return errors.New("sampler not configured")
	}

	if m.rowRenderState == nil {
		panmixer := s.GetPanMixer()

		m.rowRenderState = &state.RowRenderState{
			RenderDetails: state.RenderDetails{
				Mix:          s.Mixer(),
				SamplerSpeed: s.GetSamplerSpeed(),
				Panmixer:     panmixer,
			},
		}
	}

	var resetMemory bool
	if myCurrentRow == 0 {
		if myCurrentOrder := m.pattern.GetCurrentOrder(); myCurrentOrder == 0 {
			resetMemory = true
		}
	}

	for ch := range m.channels {
		cs := &m.channels[ch]
		cs.AdvanceRow(state.NewChannelDataTxn[TPeriod, channel.Memory](channel.GetTargetsFromData[TPeriod]))
		if resetMemory {
			mem := cs.GetMemory()
			mem.StartOrder()
		}
	}

	// generate effects and run prestart
	nch := row.GetNumChannels()
	if actual := m.GetNumChannels(); nch > actual {
		nch = actual
	}

	for channelNum := 0; channelNum < nch; channelNum++ {
		cdata := row.GetChannel(index.Channel(channelNum))

		cs := &m.channels[channelNum]
		if err := cs.SetData(cdata); err != nil {
			return err
		}
	}

	for ch := range m.channels {
		cs := &m.channels[ch]

		if txn := cs.GetTxn(); txn != nil {
			if err := txn.CommitPreRow(m, cs, cs.SemitoneSetterFactory); err != nil {
				return err
			}
		}
	}

	if err := preMixRowTxn.Commit(); err != nil {
		return err
	}

	tickDuration := tickBaseDuration / time.Duration(m.pattern.GetTempo())

	m.rowRenderState.Duration = tickDuration
	m.rowRenderState.Samples = int(tickDuration.Seconds() * float64(s.SampleRate))
	m.rowRenderState.TicksThisRow = m.pattern.GetTicksThisRow()
	m.rowRenderState.CurrentTick = 0

	// run row processing, now that prestart has completed
	for channelNum := 0; channelNum < nch; channelNum++ {
		cs := &m.channels[channelNum]

		if err := m.processRowForChannel(cs); err != nil {
			return err
		}
	}

	return nil
}

func (m *manager[TPeriod]) processRowForChannel(cs *state.ChannelState[TPeriod, channel.Memory, channel.Data]) error {
	mem := cs.GetMemory()
	mem.TremorMem().Reset()

	if txn := cs.GetTxn(); txn != nil {
		if err := txn.CommitRow(m, cs, cs.SemitoneSetterFactory); err != nil {
			return err
		}

		if err := txn.CommitPostRow(m, cs, cs.SemitoneSetterFactory); err != nil {
			return err
		}
	}
	return nil
}
