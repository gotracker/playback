package playback

import (
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/player/state"
)

// OnTick runs the S3M tick processing
func (m *Manager) OnTick() error {
	m.premix = nil
	premix, err := m.renderTick()
	if err != nil {
		return err
	}

	m.premix = premix
	return nil
}

// GetPremixData gets the current premix data from the manager
func (m *Manager) GetPremixData() (*output.PremixData, error) {
	return m.premix, nil
}

// RenderOneRow renders the next single row from the song pattern data into a RowRender object
func (m *Manager) renderTick() (*output.PremixData, error) {
	postMixRowTxn := m.pattern.StartTransaction()
	defer func() {
		postMixRowTxn.Cancel()
		m.postMixRowTxn = nil
	}()
	m.postMixRowTxn = postMixRowTxn

	if m.rowRenderState == nil || m.rowRenderState.currentTick >= m.rowRenderState.ticksThisRow {
		if err := m.processPatternRow(); err != nil {
			return nil, err
		}
	}

	var finalData render.RowRender
	premix := &output.PremixData{
		Userdata:   &finalData,
		SamplesLen: m.rowRenderState.Samples,
	}

	if err := m.soundRenderTick(premix); err != nil {
		return nil, err
	}

	finalData.Order = int(m.pattern.GetCurrentOrder())
	finalData.Row = int(m.pattern.GetCurrentRow())
	finalData.Tick = m.rowRenderState.currentTick
	if m.rowRenderState.currentTick == 0 {
		finalData.RowText = m.getRowText()
	}

	m.rowRenderState.currentTick++
	if m.rowRenderState.currentTick >= m.rowRenderState.ticksThisRow {
		postMixRowTxn.AdvanceRow = true
	}

	if err := postMixRowTxn.Commit(); err != nil {
		return nil, err
	}
	return premix, nil
}

type rowRenderState struct {
	state.RenderDetails

	ticksThisRow int
	currentTick  int
}

func (m *Manager) soundRenderTick(premix *output.PremixData) error {
	tick := m.rowRenderState.currentTick
	var lastTick = (tick+1 == m.rowRenderState.ticksThisRow)

	for ch := range m.channels {
		cs := &m.channels[ch]
		if m.song.IsChannelEnabled(ch) {

			if err := m.processEffect(ch, cs, tick, lastTick); err != nil {
				return err
			}

			rr, err := cs.RenderRowTick(m.rowRenderState.RenderDetails, nil)
			if err != nil {
				return err
			}
			if rr != nil {
				premix.Data = append(premix.Data, rr)
			}
		}
	}

	premix.MixerVolume = m.GetMixerVolume()
	return nil
}

func (m *Manager) ensureOPL2() {
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
