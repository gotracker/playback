package machine

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/vol0optimization"
)

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doNoteAction(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
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
	if m.canPastNote() {
		pn := c.cv.Clone()
		switch c.nna {
		case note.ActionCut:
			pn.Stop()
		case note.ActionRelease:
			pn.Release()
		case note.ActionFadeout:
			pn.Release()
			pn.Fadeout()
		case note.ActionRetrigger:
			pn.Release()
			pn.Attack()

		case note.ActionContinue:
			fallthrough
		default:
			// nothing
		}

		m.addPastNote(pn.(voice.RenderVoice[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]))
	}

	if pitchPanMod, ok := c.cv.(voice.PitchPanModulator[TPanning]); ok {
		pitchPanMod.SetPitchPanNote(c.prev.Semitone.Coalesce(0))
	}

	switch na.Action {
	case note.ActionCut:
		c.cv.Stop()

	case note.ActionRelease:
		c.cv.Release()

	case note.ActionFadeout:
		c.cv.Fadeout()

	case note.ActionRetrigger:
		c.cv.Release()

		if err := c.doSetupInstrument(ch, m); err != nil {
			return err
		}

		for _, o := range c.osc {
			o.Reset()
		}

		if ampMod, ok := c.cv.(voice.AmpModulator[TGlobalVolume, TMixingVolume, TVolume]); ok {
			ampMod.SetVolumeDelta(0)
		}

		if freqMod, ok := c.cv.(voice.FreqModulator[TPeriod]); ok {
			freqMod.SetPeriodDelta(0)
		}

		c.cv.Attack()

	case note.ActionContinue:
		fallthrough
	default:
		// nothing
	}
	return nil
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doSetupInstrument(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	inst := c.target.Inst
	prevInst := c.prev.Inst
	if inst != nil {
		if prevInst != inst {
			switch inst.GetKind() {
			case instrument.KindPCM:
				d := inst.GetData().(*instrument.PCM[TMixingVolume, TVolume, TPanning])
				if err := c.doSetupPCM(ch, m, inst, d); err != nil {
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
		} else {
			c.cv.Reset()
		}
	} else {
		c.cv.Stop()
	}
	return nil
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doSetupPCM(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], inst *instrument.Instrument[TMixingVolume, TVolume, TPanning], d *instrument.PCM[TMixingVolume, TVolume, TPanning]) error {
	outputRate := m.getSampleRate()

	var (
		voiceFilter  filter.Filter
		pluginFilter filter.Filter
	)
	if factory := inst.GetFilterFactory(); factory != nil {
		voiceFilter = factory(inst.SampleRate, outputRate)
	}
	if factory := inst.GetPluginFilterFactory(); factory != nil {
		pluginFilter = factory(inst.SampleRate, outputRate)
	}

	c.cv.Setup(voice.InstrumentConfig[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]{
		SampleRate:   inst.GetSampleRate(),
		AutoVibrato:  inst.GetAutoVibrato(),
		Data:         d,
		VoiceFilter:  voiceFilter,
		PluginFilter: pluginFilter,
		Vol0Optimization: vol0optimization.Vol0OptimizationSettings{
			Enabled:     true,
			MaxTicksAt0: 3,
		},
		FadeOut:              d.FadeOut,
		PitchPan:             d.PitchPan,
		VolEnv:               d.VolEnv,
		VolEnvFinishFadesOut: d.VolEnvFinishFadesOut,
		PanEnv:               d.PanEnv,
		PitchFiltMode:        d.PitchFiltMode,
		PitchFiltEnv:         d.PitchFiltEnv,
	})

	c.cv.SetPCM(d.Sample, d.Loop, d.SustainLoop, inst.GetDefaultVolume())

	return nil
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doSetupOPL2(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], inst *instrument.Instrument[TMixingVolume, TVolume, TPanning], d *instrument.OPL2) error {
	panic("unimplemented")
	//var o component.OPL2[TPeriod, TVolume]
	//o.Setup(chip, channel, reg, baseFreq, keyModulator, defaultVolume)
	//c.cv.voicer = &o
}
