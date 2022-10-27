package channel

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/player/state"
)

type dataConverter struct{}

func (c dataConverter) Process(out *state.ChannelDataActions, data *Data, cs *State) error {
	if data == nil {
		return nil
	}

	var n note.Note = note.EmptyNote{}
	inst := cs.GetInstrument()
	prevInst := inst

	if data.HasNote() || data.HasInstrument() {
		instID := data.GetInstrument(cs.StoredSemitone)
		n = data.GetNote()
		var (
			wantRetrigger    bool
			wantRetriggerVol bool
		)
		s := cs.GetSongDataInterface()
		if instID.IsEmpty() {
			// use current
			inst = prevInst
			wantRetrigger = true
		} else if !s.IsValidInstrumentID(instID) {
			out.TargetInst.Set(nil)
			n = note.InvalidNote{}
		} else {
			var str note.Semitone
			inst, str = s.GetInstrument(instID)
			n = note.CoalesceNoteSemitone(n, str)
			if !note.IsEmpty(n) && inst == nil {
				inst = prevInst
			}
			wantRetrigger = true
			wantRetriggerVol = true
		}

		if wantRetrigger {
			out.TargetInst.Set(inst)
			out.TargetPos.Set(sampling.Pos{})
			if inst != nil {
				if wantRetriggerVol {
					out.TargetVolume.Set(inst.GetDefaultVolume())
				}
				out.NoteAction.Set(note.ActionRetrigger)
			}
		}
	}

	if note.IsInvalid(n) {
		out.TargetPeriod.Set(nil)
		out.NoteAction.Set(note.ActionCut)
	} else if note.IsStop(n) {
		out.NoteAction.Set(note.ActionCut)
	} else if note.IsRelease(n) {
		out.NoteAction.Set(note.ActionRelease)
	} else if !note.IsEmpty(n) {
		if nn, ok := n.(note.Normal); ok {
			st := note.Semitone(nn)
			out.TargetStoredSemitone.Set(st)
			out.NoteCalcST.Set(st)
		} else {
			out.NoteAction.Set(note.ActionCut)
		}
	}

	if data.HasVolume() {
		v := data.GetVolume()
		if v == volume.VolumeUseInstVol {
			if inst != nil {
				v = inst.GetDefaultVolume()
			}
		}
		out.TargetVolume.Set(v)
	}

	return nil
}

type dataTxn struct {
	state.ChannelDataTxnHelper[Data, State, dataConverter]
	EffectFactory playback.EffectFactory[Data, State]
}

func (d *dataTxn) CommitPreRow(p playback.Playback, cs *State) error {
	e := d.EffectFactory(cs, d.Data)
	cs.SetActiveEffect(e)
	if e != nil {
		if onEff := p.GetOnEffect(); onEff != nil {
			onEff(e)
		}
		if err := playback.EffectPreStart(e, cs, p); err != nil {
			return err
		}
	}

	return nil
}

func (d *dataTxn) CommitRow(p playback.Playback, cs *State) error {
	if pos, ok := d.TargetPos.Get(); ok {
		cs.SetTargetPos(pos)
	}

	if inst, ok := d.TargetInst.Get(); ok {
		cs.SetTargetInst(inst)
	}

	if period, ok := d.TargetPeriod.Get(); ok {
		cs.SetTargetPeriod(period)
	}

	if st, ok := d.TargetStoredSemitone.Get(); ok {
		cs.SetStoredSemitone(st)
	}

	if v, ok := d.TargetVolume.Get(); ok {
		cs.SetActiveVolume(v)
	}

	na, targetTick := d.NoteAction.Get()
	cs.UseTargetPeriod = targetTick
	cs.SetNotePlayTick(targetTick, na, 0)

	if st, ok := d.NoteCalcST.Get(); ok {
		d.AddNoteOp(cs.SemitoneSetterFactory(st, cs.SetTargetPeriod))
	}

	return nil
}

func init() {
	var _ dataConverter
}
