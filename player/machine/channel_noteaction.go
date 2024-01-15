package machine

import (
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/voice"
)

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoNoteAction(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], outputRate frequency.Frequency) error {
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
			pn = c.cv.Clone(true)
			pn.Release()
		case note.ActionFadeout:
			pn = c.cv.Clone(true)
			pn.Release()
			pn.Fadeout()
		case note.ActionRetrigger:
			pn = c.cv.Clone(true)
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

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doSetupInstrument(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], outputRate frequency.Frequency) error {
	inst := c.target.Inst
	prevInst := c.prev.Inst
	if inst != nil {
		if prevInst != inst {
			rc := &m.actualOutputs[ch]
			if factory := inst.GetPluginFilterFactory(); factory != nil {
				rc.PluginFilter = factory(inst.SampleRate)
				rc.PluginFilter.SetPlaybackRate(outputRate)
			} else {
				rc.PluginFilter = nil
			}

			if err := c.cv.Setup(inst, outputRate); err != nil {
				return err
			}
		}
	} else {
		c.cv.Stop()
	}
	return nil
}
