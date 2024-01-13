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

	var inst *instrument.Instrument[TMixingVolume, TVolume, TPanning]
	if d.HasInstrument() {
		// retrigger (new?) instrument with period specified by `st` (0 = previous semitone)
		i := d.GetInstrument(c.prev.Semitone.Coalesce(st))

		ii, _ := m.songData.GetInstrument(i)
		inst, _ = ii.(*instrument.Instrument[TMixingVolume, TVolume, TPanning])
		wantInstrumentDefaults = inst != nil
		changeNote.Inst.Set(inst)
	} else if st != 0 {
		// retrigger same instrument
		inst = c.target.Inst
		wantInstrumentDefaults = true
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
		p := m.ConvertToPeriod(n)
		changeNote.Period.Set(p)
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
	return nil
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) decodeInstrument(m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], changeNote *NewNoteInfo[TPeriod, TMixingVolume, TVolume, TPanning], inst *instrument.Instrument[TMixingVolume, TVolume, TPanning]) error {
	if inst == nil {
		return nil
	}

	changeNote.Vol.Set(inst.GetDefaultVolume())
	changeNote.NewNoteAction.Set(inst.GetNewNoteAction())

	switch inst.GetKind() {
	case instrument.KindPCM:
		if d, ok := inst.GetData().(*instrument.PCM[TMixingVolume, TVolume, TPanning]); ok {
			if mv, set := d.MixingVolume.Get(); set {
				changeNote.MixVol.Set(mv)
			}
			if pan, set := d.Panning.Get(); set {
				changeNote.Pan.Set(pan)
			}
		}
	case instrument.KindOPL2:
		// TODO
	}
	return nil
}