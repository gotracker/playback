package machine

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/song"
)

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) decodeNote(m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], d song.ChannelDataIntf) error {
	var changeNote NewNoteInfo[TPeriod, TMixingVolume, TVolume, TPanning]

	var n note.Note
	if d.HasNote() {
		if dn := d.GetNote(); dn.Type() != note.SpecialTypeEmpty {
			n = dn
		}
	}

	var (
		st                     note.Semitone
		na                     note.Action = note.ActionContinue
		needNoteInstIdent      bool
		wantInstrumentDefaults bool
		wantTriggerNNA         bool
	)
	if n != nil {
		switch n.Type() {
		case note.SpecialTypeNormal:
			st = note.Semitone(n.(note.Normal))
			na = note.ActionRetrigger

		case note.SpecialTypeStop:
			na = note.ActionCut

		case note.SpecialTypeRelease:
			na = note.ActionRelease

		case note.SpecialTypeStopOrRelease:
			// assume cut
			na = note.ActionCut
			needNoteInstIdent = true
		}
	}

	var inst *instrument.Instrument[TPeriod, TMixingVolume, TVolume, TPanning]
	if d.HasInstrument() {
		// retrigger (new?) instrument with period specified by `st` (0 = previous semitone)
		i := d.GetInstrument()

		var ii instrument.InstrumentIntf
		ii, st = m.songData.GetInstrument(i, c.prev.Semitone.Coalesce(st))
		inst, _ = ii.(*instrument.Instrument[TPeriod, TMixingVolume, TVolume, TPanning])
		wantInstrumentDefaults = inst != nil
		changeNote.Period.Set(c.prev.Period)
		changeNote.Inst.Set(inst)
		wantTriggerNNA = true
	} else if st != 0 && c.target.Inst != nil {
		// retrigger same instrument
		i, _ := c.target.Inst.GetID().GetIndexAndSample()
		var ii instrument.InstrumentIntf
		ii, st = m.songData.GetInstrument(i, st)
		inst, _ = ii.(*instrument.Instrument[TPeriod, TMixingVolume, TVolume, TPanning])
		wantInstrumentDefaults = true
		wantTriggerNNA = true
	}

	if inst != nil && n != nil {
		if needNoteInstIdent {
			if inst.IsReleaseNote(n) {
				na = note.ActionRelease
			} else if inst.IsStopNote(n) {
				na = note.ActionCut
			}
		}
	}

	if n != nil {
		switch n.(type) {
		case note.Normal:
			// perform remap
			n = note.Normal(st)
		}
		if p := m.ConvertToPeriod(n); !p.IsInvalid() {
			changeNote.Period.Set(p)
		}
	}

	if na != note.ActionContinue {
		changeNote.ActionTick.Set(ActionTick{Action: na, Tick: 0})
		if na == note.ActionRetrigger {
			changeNote.Pos.Set(sampling.Pos{})
		}
	}

	if wantInstrumentDefaults {
		if err := c.decodeInstrument(m, &changeNote, inst); err != nil {
			return err
		}
	}

	if d.HasVolume() {
		if dd, ok := d.(song.ChannelData[TVolume]); ok {
			if v := dd.GetVolume(); !v.IsInvalid() {
				if !v.IsUseInstrumentVol() {
					changeNote.Vol.Set(v)
				}
			}
		}
	}

	c.newNote = changeNote
	c.target.TriggerNNA = wantTriggerNNA
	return nil
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) decodeInstrument(m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], changeNote *NewNoteInfo[TPeriod, TMixingVolume, TVolume, TPanning], inst *instrument.Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) error {
	if inst == nil {
		return nil
	}

	changeNote.Vol.Set(inst.GetDefaultVolume())
	changeNote.NewNoteAction.Set(inst.GetNewNoteAction())

	switch d := inst.GetData().(type) {
	case *instrument.PCM[TMixingVolume, TVolume, TPanning]:
		if pan, set := d.Panning.Get(); set {
			changeNote.Pan.Set(pan)
		}

	case *instrument.OPL2:
		// TODO - is there anything to actually do here?
	}
	return nil
}
