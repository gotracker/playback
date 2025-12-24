package machine

import (
	"errors"

	"github.com/gotracker/opl2"

	"github.com/gotracker/playback/mixing"
	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/player/sampler"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/mixer"
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

	m.hardwareSynths = append(m.hardwareSynths, opl2Synth{
		chip: m.opl2,
		gv:   func() volume.Volume { return m.gv.ToVolume() },
	})

	return nil
}

type opl2Synth struct {
	chip *opl2.Chip
	gv   func() volume.Volume
}

func (o opl2Synth) RenderTick(centerAheadPan panning.PanMixer, details mixer.Details) (mixing.Data, mixerVolumeAdjuster, error) {
	data := details.Mix.NewMixBuffer(details.Samples)
	opl2data := make([]int32, details.Samples)

	if chip := o.chip; chip != nil {
		chip.GenerateBlock2(uint(details.Samples), opl2data)
	}

	for i, s := range opl2data {
		sv := volume.Volume(s) / 32768.0
		data[i].Assign(1, []volume.Volume{sv})
	}

	mixerData := mixing.Data{
		Data:       data,
		PanMatrix:  centerAheadPan,
		Volume:     o.gv(),
		SamplesLen: details.Samples,
	}

	adjust := func(mv volume.Volume) volume.Volume {
		return mv / (mv + 1)
	}

	return mixerData, adjust, nil
}
