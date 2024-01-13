package machine

import (
	"errors"

	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/player/sampler"
)

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) setupOPL2(s *sampler.Sampler) error {
	if m.opl2 != nil || !m.ms.OPL2Enabled {
		return nil
	}

	if s != nil {
		return errors.New("sampler is nil")
	}

	var oplLen int
	//oplLen += len(m.chOrder[int(s3mfile.ChannelCategoryOPL2Melody)-1])
	//oplLen += len(m.chOrder[int(s3mfile.ChannelCategoryOPL2Drums)-1])

	if oplLen > 0 {
		o := render.NewOPL2Chip(uint32(s.SampleRate))
		o.WriteReg(0x01, 0x20) // enable all waveforms
		o.WriteReg(0x04, 0x00) // clear timer flags
		o.WriteReg(0x08, 0x40) // clear CSW and set NOTE-SEL
		o.WriteReg(0xBD, 0x00) // set default notes
		m.opl2 = o
	}

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) renderOPL2Tick(mixerData *mixing.Data, mix *mixing.Mixer, tickSamples int) error {
	// make a stand-alone data buffer for this channel for this tick
	data := mix.NewMixBuffer(tickSamples)

	opl2data := make([]int32, tickSamples)

	if opl2 := m.opl2; opl2 != nil {
		opl2.GenerateBlock2(uint(tickSamples), opl2data)
	}

	for i, s := range opl2data {
		sv := volume.Volume(s) / 32768.0
		data[i].Assign(1, []volume.Volume{sv})
	}
	*mixerData = mixing.Data{
		Data:       data,
		Pan:        panning.CenterAhead,
		Volume:     m.gv.ToVolume(),
		SamplesLen: tickSamples,
	}
	return nil
}
