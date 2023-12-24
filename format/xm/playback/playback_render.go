package playback

import (
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/player/render"
)

// OnTick runs the XM tick processing
func (m *manager[TPeriod]) OnTick() error {
	m.premix = nil

	postMixRowTxn := m.pattern.StartTransaction()
	defer func() {
		postMixRowTxn.Cancel()
		m.postMixRowTxn = nil
	}()
	m.postMixRowTxn = postMixRowTxn

	if m.rowRenderState == nil || m.rowRenderState.CurrentTick >= m.rowRenderState.TicksThisRow {
		if err := m.processPatternRow(); err != nil {
			return err
		}
	}

	var finalData render.RowRender
	premix := &output.PremixData{
		Userdata:   &finalData,
		SamplesLen: m.rowRenderState.Samples,
	}

	if err := m.soundRenderTick(premix); err != nil {
		return err
	}

	finalData.Order = int(m.pattern.GetCurrentOrder())
	finalData.Row = int(m.pattern.GetCurrentRow())
	finalData.Tick = m.rowRenderState.CurrentTick
	if m.rowRenderState.CurrentTick == 0 {
		finalData.RowText = m.getRowText()
	}

	m.rowRenderState.CurrentTick++
	if m.rowRenderState.CurrentTick >= m.rowRenderState.TicksThisRow {
		postMixRowTxn.AdvanceRow = true
	}

	if err := postMixRowTxn.Commit(); err != nil {
		return err
	}

	m.premix = premix
	return nil
}

// GetPremixData gets the current premix data from the manager
func (m *manager[TPeriod]) GetPremixData() (*output.PremixData, error) {
	return m.premix, nil
}

func (m *manager[TPeriod]) soundRenderTick(premix *output.PremixData) error {
	tick := m.rowRenderState.CurrentTick
	var lastTick = (tick+1 == m.rowRenderState.TicksThisRow)

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

/** unused in XM, for now
func (m *manager[TPeriod]) ensureOPL2() {
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
*/
