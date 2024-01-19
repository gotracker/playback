package machine

import (
	"errors"
	"fmt"

	"github.com/gotracker/gomixing/sampling"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/oscillator"
	"github.com/gotracker/playback/voice/types"
)

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) canPastNote() bool {
	return m.us.EnableNewNoteActions && len(m.virtualOutputs) > 0
}

func withChannel[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, fn func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error) error {
	c, err := m.getChannel(ch)
	if err != nil {
		return err
	}

	return fn(c)
}

func withChannelReturningValue[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning, T any](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, fn func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) (T, error)) (T, error) {
	c, err := m.getChannel(ch)
	if err != nil {
		var empty T
		return empty, err
	}

	return fn(c)
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetChannelMemory(ch index.Channel) (song.ChannelMemory, error) {
	return withChannelReturningValue(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) (song.ChannelMemory, error) {
		return c.memory, nil
	})
}

func GetChannelMemory[TMemory song.ChannelMemory, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel) (TMemory, error) {
	chMem, err := m.GetChannelMemory(ch)
	if err != nil {
		var empty TMemory
		return empty, err
	}

	mem, ok := chMem.(TMemory)
	if !ok {
		var empty TMemory
		return empty, errors.New("could not convert channel memory type")
	}

	return mem, nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) IsChannelMuted(ch index.Channel) (bool, error) {
	return withChannelReturningValue(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) (bool, error) {
		return c.cv.IsMuted(), nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelMute(ch index.Channel, muted bool) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		prevMute := c.cv.IsMuted()
		traceChannelValueChangeWithComment(m, ch, "mute", prevMute, muted, "SetChannelMute")
		return c.cv.SetMuted(muted)
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelMixingVolume(ch index.Channel, v TMixingVolume) error {
	if v.IsInvalid() {
		return fmt.Errorf("channel[%d] mixing volume out of range: %v", ch, v)
	}

	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if ampMod, ok := c.cv.(voice.AmpModulator[TGlobalVolume, TMixingVolume, TVolume]); ok {
			mv := ampMod.GetMixingVolume()
			traceChannelValueChangeWithComment(m, ch, "mv", mv, v, "SetChannelMixingVolume")
			ampMod.SetMixingVolume(v)
		}
		return nil
	})
}

func (m machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetChannelPeriod(ch index.Channel) (TPeriod, error) {
	return withChannelReturningValue(&m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) (TPeriod, error) {
		if freqMod, ok := c.cv.(voice.FreqModulator[TPeriod]); ok {
			return freqMod.GetPeriod(), nil
		}
		var empty TPeriod
		return empty, nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelPeriod(ch index.Channel, p TPeriod) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if freqMod, ok := c.cv.(voice.FreqModulator[TPeriod]); ok {
			tp := freqMod.GetPeriod()
			traceChannelValueChangeWithComment(m, ch, "target.Period", tp, p, "SetChannelPeriod")
			freqMod.SetPeriod(p)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelPeriodDelta(ch index.Channel, d period.Delta) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if freqMod, ok := c.cv.(voice.FreqModulator[TPeriod]); ok {
			pd := freqMod.GetPeriodDelta()
			traceChannelValueChangeWithComment(m, ch, "pd", pd, d, "SetChannelPeriodDelta")
			freqMod.SetPeriodDelta(d)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetChannelInstrument(ch index.Channel) (*instrument.Instrument[TPeriod, TMixingVolume, TVolume, TPanning], error) {
	return withChannelReturningValue(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) (*instrument.Instrument[TPeriod, TMixingVolume, TVolume, TPanning], error) {
		return c.target.Inst, nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelInstrumentByID(ch index.Channel, i instrument.ID) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		var oldID instrument.ID
		if c.target.Inst != nil {
			oldID = c.target.Inst.GetID()
		}

		traceChannelValueChangeWithComment(m, ch, "target.Inst", oldID, i, "SetChannelInstrumentByID")
		inst, _ := m.songData.GetInstrument(i)
		var ok bool
		c.target.Inst, ok = inst.(*instrument.Instrument[TPeriod, TMixingVolume, TVolume, TPanning])
		if !ok {
			return errors.New("could not convert instrument to pointer type")
		}

		return nil
	})
}

func (m machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetChannelVolume(ch index.Channel) (TVolume, error) {
	return withChannelReturningValue(&m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) (TVolume, error) {
		if ampMod, ok := c.cv.(voice.AmpModulator[TGlobalVolume, TMixingVolume, TVolume]); ok {
			return ampMod.GetVolume(), nil
		}
		var empty TVolume
		return empty, nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelVolume(ch index.Channel, v TVolume) error {
	if v.IsInvalid() {
		return fmt.Errorf("channel[%d] volume out of range: %v", ch, v)
	}

	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		traceChannelOptionalValueChangeWithComment(m, ch, "newNote.Vol", c.newNote.Vol, v, "SetChannelVolume")
		c.newNote.Vol.Set(v)
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelVolumeDelta(ch index.Channel, d types.VolumeDelta) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if ampMod, ok := c.cv.(voice.AmpModulator[TGlobalVolume, TMixingVolume, TVolume]); ok {
			vd := ampMod.GetVolumeDelta()
			traceChannelValueChangeWithComment(m, ch, "vd", vd, d, "SetChannelVolumeDelta")
			ampMod.SetVolumeDelta(d)
		}
		return nil
	})
}

func (m machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetChannelPan(ch index.Channel) (TPanning, error) {
	return withChannelReturningValue(&m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) (TPanning, error) {
		if panMod, ok := c.cv.(voice.PanModulator[TPanning]); ok {
			return panMod.GetPan(), nil
		}
		return types.GetPanDefault[TPanning](), nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelPan(ch index.Channel, pan TPanning) error {
	if pan.IsInvalid() {
		return fmt.Errorf("channel[%d] channel pan out of range: %v", ch, pan)
	}

	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if panMod, ok := c.cv.(voice.PanModulator[TPanning]); ok {
			orig := panMod.GetPan()
			traceChannelValueChangeWithComment(m, ch, "pan", orig, pan, "SetChannelPan")
			panMod.SetPan(pan)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelPanningDelta(ch index.Channel, d types.PanDelta) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if panMod, ok := c.cv.(voice.PanModulator[TPanning]); ok {
			orig := panMod.GetPanDelta()
			traceChannelValueChangeWithComment(m, ch, "delta", orig, d, "SetChannelPanningDelta")
			panMod.SetPanDelta(d)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelSurround(ch index.Channel, enabled bool) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		traceChannelValueChangeWithComment(m, ch, "surround", c.surround, enabled, "SetChannelSurround")
		c.surround = enabled
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelFilter(ch index.Channel, f filter.Filter) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		traceChannelValueChangeWithComment(m, ch, "filter", c.filter, f, "SetChannelFilter")
		c.filter = f
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) ChannelStopOrRelease(ch index.Channel) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		traceChannel(m, ch, "ChannelStopOrRelease")
		var n note.StopOrReleaseNote
		if c.target.Inst != nil && c.target.Inst.IsReleaseNote(n) {
			c.cv.Release()
			return nil
		}

		c.cv.Stop()
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) ChannelStop(ch index.Channel) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		traceChannel(m, ch, "ChannelStop")
		c.cv.Stop()
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) ChannelRelease(ch index.Channel) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		traceChannel(m, ch, "ChannelRelease")
		c.cv.Release()
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) ChannelFadeout(ch index.Channel) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		traceChannel(m, ch, "ChannelFadeout")
		c.cv.Fadeout()
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetNextChannelWavetableValue(ch index.Channel, speed int, depth float32, oscSelect Oscillator) (float32, error) {
	if int(oscSelect) >= NumOscillators {
		return 0, fmt.Errorf("oscillator select out of range: %v", oscSelect)
	}

	return withChannelReturningValue(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) (float32, error) {
		o := c.osc[oscSelect]
		out := o.GetWave(depth)
		o.Advance(speed)
		return out, nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelNoteAction(ch index.Channel, na note.Action, tick int) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		at := ActionTick{Action: na, Tick: tick}
		traceChannelOptionalValueChangeWithComment(m, ch, "target.ActionTick", c.target.ActionTick, at, "SetChannelNoteAction")
		c.target.ActionTick.Set(at)
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetPatternLoopStart(ch index.Channel) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		traceChannelValueChangeWithComment(m, ch, "patternLoopStart", c.patternLoop.Start, m.ticker.current.Row, "SetPatternLoopStart")
		c.patternLoop.Start = m.ticker.current.Row
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetPatternLoops(ch index.Channel, count int) error {
	if count <= 0 {
		return fmt.Errorf("loop count out of range: %d", count)
	}

	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		pl := &c.patternLoop

		traceChannelValueChangeWithComment(m, ch, "patternLoopEnd", pl.End, m.ticker.current.Row, "SetPatternLoops")
		pl.End = m.ticker.current.Row

		traceChannelValueChangeWithComment(m, ch, "patternLoopTotal", pl.Total, count, "SetPatternLoops")
		pl.Total = count

		return nil
	})

}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) StartChannelPortaToNote(ch index.Channel) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		traceChannelOptionalValueResetWithComment(m, ch, "newNote.ActionTick", c.newNote.ActionTick, "StartChannelPortaToNote")
		c.newNote.ActionTick.Reset()

		if p, set := c.newNote.Period.Get(); set {
			traceChannelOptionalValueResetWithComment(m, ch, "newNote.Period", c.newNote.Period, "StartChannelPortaToNote")
			c.newNote.Period.Reset()
			traceChannelValueChangeWithComment(m, ch, "target.PortaPeriod", c.target.PortaPeriod, p, "StartChannelPortaToNote")
			c.target.PortaPeriod = p
		}

		traceChannelOptionalValueResetWithComment(m, ch, "newNote.Pos", c.newNote.Pos, "StartChannelPortaToNote")
		c.newNote.Pos.Reset()
		c.target.TriggerNNA = false

		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoChannelPortaToNote(ch index.Channel, delta period.Delta) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if freqMod, ok := c.cv.(voice.FreqModulator[TPeriod]); ok {
			var p TPeriod
			if m.ms.Quirks.PortaToNoteUsesModifiedPeriod {
				var err error
				p, err = freqMod.GetFinalPeriod()
				if err != nil {
					return err
				}
			} else {
				p = freqMod.GetPeriod()
			}
			tp, err := m.ms.PeriodConverter.PortaToNote(p, delta, c.target.PortaPeriod)
			if err != nil {
				return err
			}

			traceChannelValueChangeWithComment(m, ch, "target.Period", p, tp, "DoChannelPortaToNote (%d)", delta)
			freqMod.SetPeriod(tp)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoChannelPortaDown(ch index.Channel, delta period.Delta) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if freqMod, ok := c.cv.(voice.FreqModulator[TPeriod]); ok {
			p := freqMod.GetPeriod()
			tp, err := m.ms.PeriodConverter.PortaDown(p, delta)
			if err != nil {
				return err
			}

			traceChannelValueChangeWithComment(m, ch, "target.Period", p, tp, "DoChannelPortaDown (%d)", delta)
			freqMod.SetPeriod(tp)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoChannelPortaUp(ch index.Channel, delta period.Delta) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if freqMod, ok := c.cv.(voice.FreqModulator[TPeriod]); ok {
			p := freqMod.GetPeriod()
			tp, err := m.ms.PeriodConverter.PortaUp(p, delta)
			if err != nil {
				return err
			}

			traceChannelValueChangeWithComment(m, ch, "target.Period", p, tp, "DoChannelPortaUp (%d)", delta)
			freqMod.SetPeriod(tp)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoChannelArpeggio(ch index.Channel, delta int8) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if freqMod, ok := c.cv.(voice.FreqModulator[TPeriod]); ok {
			p := freqMod.GetPeriod()

			st := c.prev.Semitone.Coalesce(0)

			tp := m.ConvertToPeriod(note.Normal(note.Semitone(int8(st) + delta)))

			traceChannelValueChangeWithComment(m, ch, "target.Period", p, tp, "DoChannelArpeggio")
			freqMod.SetPeriod(tp)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SlideChannelVolume(ch index.Channel, multiplier, add float32) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if ampMod, ok := c.cv.(voice.AmpModulator[TGlobalVolume, TMixingVolume, TVolume]); ok {
			vol := ampMod.GetVolume()
			fma, ok := any(vol).(VolumeFMA[TVolume])
			if !ok {
				return errors.New("could not determine FMA interface for channel volume")
			}
			v := fma.FMA(multiplier, add)

			if v.IsInvalid() {
				return fmt.Errorf("channel volume out of range: %v", v)
			}

			traceChannelValueChangeWithComment[TVolume](m, ch, "vol", vol, v, "SlideChannelVolume")
			ampMod.SetVolume(v)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SlideChannelMixingVolume(ch index.Channel, multiplier, add float32) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if ampMod, ok := c.cv.(voice.AmpModulator[TGlobalVolume, TMixingVolume, TVolume]); ok {
			mv := ampMod.GetMixingVolume()
			fma, ok := any(mv).(VolumeFMA[TMixingVolume])
			if !ok {
				return errors.New("could not determine FMA interface for channel mixing volume")
			}
			v := fma.FMA(multiplier, add)

			if v.IsInvalid() {
				return fmt.Errorf("channel mixing volume out of range: %v", v)
			}

			traceChannelValueChangeWithComment[TMixingVolume](m, ch, "mv", mv, v, "SlideChannelMixingVolume")
			ampMod.SetMixingVolume(v)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelPos(ch index.Channel, pos sampling.Pos) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		traceChannelOptionalValueChangeWithComment(m, ch, "newNote.Pos", c.newNote.Pos, pos, "SetChannelPos")
		c.newNote.Pos.Set(pos)
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelEnvelopePositions(ch index.Channel, pos int) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if volEnv, ok := c.cv.(voice.VolumeEnvelope[TGlobalVolume, TMixingVolume, TVolume]); ok {
			volPos := volEnv.GetVolumeEnvelopePosition()
			traceChannelValueChangeWithComment(m, ch, "volEnv.Pos", volPos, pos, "SetChannelEnvelopePositions")
			volEnv.SetVolumeEnvelopePosition(pos)
		}

		if pitchEnv, ok := c.cv.(voice.PitchEnvelope[TPeriod]); ok {
			pitchPos := pitchEnv.GetPitchEnvelopePosition()
			traceChannelValueChangeWithComment(m, ch, "pitchEnv.Pos", pitchPos, pos, "SetChannelEnvelopePositions")
			pitchEnv.SetPitchEnvelopePosition(pos)
		}

		if panEnv, ok := c.cv.(voice.PanEnvelope[TPanning]); ok {
			panPos := panEnv.GetPanEnvelopePosition()
			traceChannelValueChangeWithComment(m, ch, "panEnv.Pos", panPos, pos, "SetChannelEnvelopePositions")
			panEnv.SetPanEnvelopePosition(pos)
		}

		if filterEnv, ok := c.cv.(voice.FilterEnvelope); ok {
			filtPos := filterEnv.GetFilterEnvelopePosition()
			traceChannelValueChangeWithComment(m, ch, "filterEnv.Pos", filtPos, pos, "SetChannelEnvelopePositions")
			filterEnv.SetFilterEnvelopePosition(pos)
		}

		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SlideChannelPan(ch index.Channel, multiplier, add float32) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if panMod, ok := c.cv.(voice.PanModulator[TPanning]); ok {
			p := panMod.GetPan()
			fma, ok := any(p).(PanningFMA[TPanning])
			if !ok {
				return errors.New("could not determine FMA interface for channel panning")
			}
			v := fma.FMA(multiplier, add)

			if v.IsInvalid() {
				return fmt.Errorf("channel panning out of range: %v", v)
			}

			traceChannelValueChangeWithComment[TPanning](m, ch, "pan", p, v, "SlideChannelPan")
			panMod.SetPan(v)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelVolumeActive(ch index.Channel, on bool) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if ampMod, ok := c.cv.(voice.AmpModulator[TGlobalVolume, TMixingVolume, TVolume]); ok {
			active := ampMod.IsActive()
			traceChannelValueChangeWithComment(m, ch, "active", active, on, "SetChannelVolumeActive")
			ampMod.SetActive(on)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelOscillatorWaveform(ch index.Channel, osc Oscillator, wave oscillator.WaveTableSelect) error {
	if int(osc) >= NumOscillators {
		return fmt.Errorf("oscillator select out of range: %v", osc)
	}

	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		wf := c.osc[osc].GetWaveform()
		traceChannelValueChangeWithComment(m, ch, fmt.Sprintf("osc[%d].wave", osc), wf, wave, "SetChannelOscillatorWaveform")
		c.osc[osc].SetWaveform(wave)
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoChannelPastNoteEffect(ch index.Channel, na note.Action) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		c.doPastNoteAction(m, na)
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelNewNoteAction(ch index.Channel, na note.Action) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		traceChannelValueChangeWithComment(m, ch, "nna", c.nna, na, "SetChannelNewNoteAction")
		c.nna = na
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelVolumeEnvelopeEnable(ch index.Channel, enabled bool) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if volEnv, ok := c.cv.(voice.VolumeEnvelope[TGlobalVolume, TMixingVolume, TVolume]); ok {
			on := volEnv.IsVolumeEnvelopeEnabled()
			traceChannelValueChangeWithComment(m, ch, "volEnv.enabled", on, enabled, "SetChannelVolumeEnvelopeEnable")
			volEnv.EnableVolumeEnvelope(enabled)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelPanningEnvelopeEnable(ch index.Channel, enabled bool) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if panEnv, ok := c.cv.(voice.PanEnvelope[TPanning]); ok {
			on := panEnv.IsPanEnvelopeEnabled()
			traceChannelValueChangeWithComment(m, ch, "panEnv.enabled", on, enabled, "SetChannelPanningEnvelopeEnable")
			panEnv.EnablePanEnvelope(enabled)
		}
		return nil
	})
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetChannelPitchEnvelopeEnable(ch index.Channel, enabled bool) error {
	return withChannel(m, ch, func(c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
		if pitchEnv, ok := c.cv.(voice.PitchEnvelope[TPeriod]); ok {
			on := pitchEnv.IsPitchEnvelopeEnabled()
			traceChannelValueChangeWithComment(m, ch, "pitchEnv.enabled", on, enabled, "SetChannelPitchEnvelopeEnable")
			pitchEnv.EnablePitchEnvelope(enabled)
		}
		return nil
	})
}
