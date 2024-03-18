package machine

import (
	"errors"

	"github.com/gotracker/opl2"
	"github.com/gotracker/playback/mixing"
	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/player/sampler"
	"github.com/gotracker/playback/voice"
)

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) setupOPL2(s *sampler.Sampler) error {
	if s == nil {
		return errors.New("sampler is nil")
	}

	o := opl2.NewChip(uint32(s.SampleRate), false)
	o.WriteReg(0x01, 0x20) // enable all waveforms
	o.WriteReg(0x04, 0x00) // clear timer flags
	o.WriteReg(0x08, 0x40) // clear CSW and set NOTE-SEL
	o.WriteReg(0xBD, 0x00) // set default notes
	m.opl2 = o

	for i := range m.actualOutputs {
		rc := &m.actualOutputs[i]
		if v, _ := rc.GetVoice().(voice.VoiceOPL2er); v != nil {
			v.SetOPL2Chip(m.opl2)
		}
	}

	for i := range m.virtualOutputs {
		rc := &m.virtualOutputs[i]
		if v, _ := rc.GetVoice().(voice.VoiceOPL2er); v != nil {
			v.SetOPL2Chip(m.opl2)
		}
	}

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) renderOPL2Tick(centerAheadPan panning.PanMixer, mixerData *mixing.Data, mix *mixing.Mixer, tickSamples int) error {
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
		PanMatrix:  centerAheadPan,
		Volume:     m.gv.ToVolume(),
		SamplesLen: tickSamples,
	}
	return nil
}
