package machine

import (
	"fmt"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
)

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) OrderStart(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	if m.ticker.current.Order == 0 {
		c.memory.StartOrder0()
	}
	c.resetPatternLoop()
	return nil
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) RowStart(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	traceChannelOptionalValueResetWithComment(m, ch, "target.ActionTick", c.target.ActionTick, "channel.RowStart")
	c.target.ActionTick.Reset()

	if freqMod, ok := c.cv.(voice.FreqModulator[TPeriod]); ok {
		pd := freqMod.GetPeriodDelta()
		var reset period.Delta
		traceChannelValueChangeWithComment(m, ch, "pd", pd, reset, "channel.RowStart")
		freqMod.SetPeriodDelta(reset)
	}

	for _, i := range c.instructions {
		if err := m.DoInstructionRowStart(ch, i); err != nil {
			return err
		}
	}

	info := c.newNote
	c.newNote.Reset()

	if tp, set := info.Period.Get(); set {
		if freqMod, ok := c.cv.(voice.FreqModulator[TPeriod]); ok {
			orig := freqMod.GetPeriod()
			traceChannelValueChangeWithComment(m, ch, "target.Period", orig, tp, "channel.RowStart")
			freqMod.SetPeriod(tp)
		}
	}

	if inst, set := info.Inst.Get(); set {
		var prev, next instrument.ID
		if c.target.Inst != nil {
			prev = c.target.Inst.GetID()
		}
		if inst != nil {
			next = inst.GetID()
		}

		traceChannelValueChangeWithComment(m, ch, "target.Inst", prev, next, "channel.RowStart")
		c.target.Inst = inst
	}

	if tpos, set := info.Pos.Get(); set {
		traceChannelOptionalValueChangeWithComment(m, ch, "target.Pos", c.target.Pos, tpos, "channel.RowStart")
		c.target.Pos.Set(tpos)
	}

	if ampMod, ok := c.cv.(voice.AmpModulator[TGlobalVolume, TMixingVolume, TVolume]); ok {
		if mv, set := info.MixVol.Get(); set {
			if mv.IsInvalid() {
				return fmt.Errorf("channel[%d] mixing volume out of range: %v", ch, mv)
			}

			orig := ampMod.GetMixingVolume()
			traceChannelValueChangeWithComment(m, ch, "mv", orig, mv, "channel.RowStart")
			ampMod.SetMixingVolume(mv)
		}

		if vol, set := info.Vol.Get(); set {
			if vol.IsInvalid() {
				return fmt.Errorf("channel[%d] volume out of range: %v", ch, vol)
			}

			orig := ampMod.GetVolume()
			traceChannelValueChangeWithComment(m, ch, "vol", orig, vol, "channel.RowStart")
			ampMod.SetVolume(vol)
		}
	}

	if pan, set := info.Pan.Get(); set {
		if pan.IsInvalid() {
			return fmt.Errorf("channel[%d] channel pan out of range: %v", ch, pan)
		}

		if panMod, ok := c.cv.(voice.PanModulator[TPanning]); ok {
			orig := panMod.GetPan()
			traceChannelValueChangeWithComment(m, ch, "pan", orig, pan, "channel.RowStart")
			panMod.SetPan(pan)
		}
	}

	if nna, set := info.NewNoteAction.Get(); set {
		traceChannelValueChangeWithComment(m, ch, "nna", c.nna, nna, "channel.RowStart")
		c.nna = nna
	}

	if na, set := info.ActionTick.Get(); set {
		traceChannelOptionalValueChangeWithComment(m, ch, "target.ActionTick", c.target.ActionTick, na, "channel.RowStart")
		c.target.ActionTick.Set(na)
	}

	if c.target.Pos.IsSet() && !c.target.ActionTick.IsSet() {
		c.target.ActionTick.Set(ActionTick{
			Action: note.ActionRetrigger,
			Tick:   0,
		})
	}

	return nil
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) Tick(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	for _, i := range c.instructions {
		if err := m.DoInstructionTick(ch, i); err != nil {
			return err
		}
	}

	return nil
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) RowEnd(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	for _, i := range c.instructions {
		if err := m.DoInstructionRowEnd(ch, i); err != nil {
			return err
		}
	}

	var (
		prevID instrument.ID
		curID  instrument.ID
	)
	if c.prev.Inst != nil {
		prevID = c.prev.Inst.GetID()
	}
	if c.target.Inst != nil {
		curID = c.target.Inst.GetID()
	}
	traceChannelValueChangeWithComment(m, ch, "prev.Inst", prevID, curID, "channel.RowEnd")
	c.prev.Inst = c.target.Inst

	if freqMod, ok := c.cv.(voice.FreqModulator[TPeriod]); ok {
		var p TPeriod
		if m.ms.Quirks.PreviousPeriodUsesModifiedPeriod {
			var err error
			p, err = freqMod.GetFinalPeriod()
			if err != nil {
				return err
			}
		} else {
			p = freqMod.GetPeriod()
		}
		traceChannelValueChangeWithComment(m, ch, "prev.Period", c.prev.Period, p, "channel.RowEnd")
		c.prev.Period = p
	}

	if err := c.doPatternLoop(ch, m); err != nil {
		return nil
	}

	return nil
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) OrderEnd(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	return nil
}
