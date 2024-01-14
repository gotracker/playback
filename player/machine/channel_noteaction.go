package machine

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/system"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/types"
)

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoNoteAction(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], outputRate system.Frequency) error {
	na, set := c.target.ActionTick.Get()
	if !set {
		// assume continue
		return nil
	}

	if na.Tick != m.ticker.current.tick {
		// not time yet
		return nil
	}

	// consume the action
	traceChannelOptionalValueResetWithComment(m, ch, "target.ActionTick", c.target.ActionTick, "doNoteAction")
	c.target.ActionTick.Reset()

	// perform new note action
	if na.Action != note.ActionContinue && m.canPastNote() {
		var pn voice.Voice
		switch c.nna {
		case note.ActionCut:
			c.cv.Stop()
		case note.ActionRelease:
			pn = c.cv.Clone()
			pn.Release()
		case note.ActionFadeout:
			pn = c.cv.Clone()
			pn.Release()
			pn.Fadeout()
		case note.ActionRetrigger:
			pn = c.cv.Clone()
			pn.Release()
			pn.Attack()

		case note.ActionContinue:
			fallthrough
		default:
			// nothing
		}

		if pn != nil {
			c.addPastNote(m, pn.(voice.RenderVoice[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]))
		}
	}

	switch na.Action {
	case note.ActionCut:
		c.cv.Stop()
		return nil

	case note.ActionRelease:
		c.cv.Release()

	case note.ActionFadeout:
		c.cv.Fadeout()

	case note.ActionRetrigger:
		c.cv.Release()

		if err := c.doSetupInstrument(ch, m, outputRate); err != nil {
			return err
		}

		c.memory.Retrigger()

		for _, o := range c.osc {
			o.Reset()
		}

		c.cv.Reset()

		c.cv.Attack()

	case note.ActionContinue:
		fallthrough
	default:
		// nothing
	}

	if pitchPanMod, ok := c.cv.(voice.PitchPanModulator[TPanning]); ok {
		pitchPanMod.SetPitchPanNote(c.prev.Semitone.Coalesce(0))
	}

	if pos, set := c.target.Pos.Get(); set {
		if samp, ok := c.cv.(voice.Sampler); ok {
			samp.SetPos(pos)
		}
		c.target.Pos.Reset()
	}

	return nil
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doSetupInstrument(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], outputRate system.Frequency) error {
	inst := c.target.Inst
	prevInst := c.prev.Inst
	if inst != nil {
		if prevInst != inst {
			switch inst.GetKind() {
			case instrument.KindPCM:
				d := inst.GetData().(*instrument.PCM[TMixingVolume, TVolume, TPanning])
				if err := c.doSetupPCM(ch, m, inst, d, outputRate); err != nil {
					return err
				}

			case instrument.KindOPL2:
				d := inst.GetData().(*instrument.OPL2)
				if err := c.doSetupOPL2(ch, m, inst, d); err != nil {
					return err
				}

			default:
				panic("unhandled instrument kind")
			}
		}
	} else {
		c.cv.Stop()
	}
	return nil
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doSetupPCM(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], inst *instrument.Instrument[TMixingVolume, TVolume, TPanning], d *instrument.PCM[TMixingVolume, TVolume, TPanning], outputRate system.Frequency) error {
	var voiceFilter filter.Filter
	if factory := inst.GetFilterFactory(); factory != nil {
		voiceFilter = factory(inst.SampleRate)
		voiceFilter.SetPlaybackRate(outputRate)
	}

	rc := &m.actualOutputs[ch]
	if factory := inst.GetPluginFilterFactory(); factory != nil {
		rc.PluginFilter = factory(inst.SampleRate)
		rc.PluginFilter.SetPlaybackRate(outputRate)
	} else {
		rc.PluginFilter = nil
	}

	c.cv.Setup(voice.InstrumentConfig[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]{
		SampleRate:           inst.GetSampleRate(),
		AutoVibrato:          inst.GetAutoVibrato(),
		Data:                 d,
		VoiceFilter:          voiceFilter,
		FadeOut:              d.FadeOut,
		PitchPan:             d.PitchPan,
		VolEnv:               d.VolEnv,
		VolEnvFinishFadesOut: d.VolEnvFinishFadesOut,
		PanEnv:               d.PanEnv,
		PitchFiltMode:        d.PitchFiltMode,
		PitchFiltEnv:         d.PitchFiltEnv,
	})

	var mixVol TMixingVolume
	if mv, set := d.MixingVolume.Get(); set {
		mixVol = mv
	} else {
		mixVol = types.GetMaxVolume[TMixingVolume]()
	}

	c.cv.SetPCM(d.Sample, d.Loop, d.SustainLoop, mixVol, inst.GetDefaultVolume())

	return nil
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doSetupOPL2(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], inst *instrument.Instrument[TMixingVolume, TVolume, TPanning], d *instrument.OPL2) error {
	panic("unimplemented")
	//var o component.OPL2[TPeriod, TVolume]
	//o.Setup(chip, channel, reg, baseFreq, keyModulator, defaultVolume)
	//c.cv.voicer = &o
}
