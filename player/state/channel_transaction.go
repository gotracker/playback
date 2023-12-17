package state

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
	"github.com/heucuva/optional"
)

type ChannelDataTransaction[TMemory any] interface {
	GetData() song.ChannelData
	SetData(data song.ChannelData, s song.Data, cs *ChannelState[TMemory]) error

	CommitPreRow(p playback.Playback, cs *ChannelState[TMemory], semitoneSetterFactory SemitoneSetterFactory[TMemory]) error
	CommitRow(p playback.Playback, cs *ChannelState[TMemory], semitoneSetterFactory SemitoneSetterFactory[TMemory]) error
	CommitPostRow(p playback.Playback, cs *ChannelState[TMemory], semitoneSetterFactory SemitoneSetterFactory[TMemory]) error

	CommitPreTick(p playback.Playback, cs *ChannelState[TMemory], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory]) error
	CommitTick(p playback.Playback, cs *ChannelState[TMemory], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory]) error
	CommitPostTick(p playback.Playback, cs *ChannelState[TMemory], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory]) error

	AddVolOp(op VolOp[TMemory])
	AddNoteOp(op NoteOp[TMemory])
}

type ChannelDataActions struct {
	NoteAction optional.Value[note.Action]
	NoteCalcST optional.Value[note.Semitone]

	TargetPos            optional.Value[sampling.Pos]
	TargetInst           optional.Value[*instrument.Instrument]
	TargetPeriod         optional.Value[period.Period]
	TargetStoredSemitone optional.Value[note.Semitone]
	TargetNewNoteAction  optional.Value[note.Action]
	TargetVolume         optional.Value[volume.Volume]
}

type ChannelDataConverter[TMemory any] func(out *ChannelDataActions, data song.ChannelData, s song.Data, cs *ChannelState[TMemory]) error

type channelDataTxnHelper[TMemory any] struct {
	Data          song.ChannelData
	effectFactory func(*TMemory, song.ChannelData) playback.Effect

	ChannelDataActions

	VolOps  []VolOp[TMemory]
	NoteOps []NoteOp[TMemory]
}

func NewChannelDataTxn[TMemory any](effectFactory func(*TMemory, song.ChannelData) playback.Effect) ChannelDataTransaction[TMemory] {
	return &channelDataTxnHelper[TMemory]{
		effectFactory: effectFactory,
	}
}

func (d *channelDataTxnHelper[TMemory]) GetData() song.ChannelData {
	return d.Data
}

func (d *channelDataTxnHelper[TMemory]) SetData(cd song.ChannelData, s song.Data, cs *ChannelState[TMemory]) error {
	d.Data = cd

	if d.ProcessData != nil {
		return d.ProcessData(&d.ChannelDataActions, cd, s, cs)
	}
	return nil
}

func (d *channelDataTxnHelper[TMemory]) CommitPreRow(p playback.Playback, cs *ChannelState[TMemory], semitoneSetterFactory SemitoneSetterFactory[TMemory]) error {
	e := d.effectFactory(cs.GetMemory(), d.Data)
	cs.SetActiveEffect(e)
	if e != nil {
		if onEff := p.GetOnEffect(); onEff != nil {
			onEff(e)
		}
		if err := playback.EffectPreStart[TMemory](e, cs, p); err != nil {
			return err
		}
	}

	return nil
}

func (d *channelDataTxnHelper[TMemory]) CommitRow(p playback.Playback, cs *ChannelState[TMemory], semitoneSetterFactory SemitoneSetterFactory[TMemory]) error {
	if pos, ok := d.TargetPos.Get(); ok {
		cs.SetTargetPos(pos)
	}

	if inst, ok := d.TargetInst.Get(); ok {
		cs.SetTargetInst(inst)
	}

	if period, ok := d.TargetPeriod.Get(); ok {
		cs.SetTargetPeriod(period)
		cs.SetPortaTargetPeriod(period)
	}

	if st, ok := d.TargetStoredSemitone.Get(); ok {
		cs.SetStoredSemitone(st)
	}

	if nna, ok := d.TargetNewNoteAction.Get(); ok {
		cs.SetNewNoteAction(nna)
	}

	if v, ok := d.TargetVolume.Get(); ok {
		cs.SetActiveVolume(v)
	}

	na, targetTick := d.NoteAction.Get()
	cs.UseTargetPeriod = targetTick
	cs.SetNotePlayTick(targetTick, na, 0)

	if st, ok := d.NoteCalcST.Get(); ok {
		d.AddNoteOp(semitoneSetterFactory(st, cs.SetTargetPeriod))
	}

	return nil
}

func (d *channelDataTxnHelper[TMemory]) CommitPostRow(p playback.Playback, cs *ChannelState[TMemory], semitoneSetterFactory SemitoneSetterFactory[TMemory]) error {
	return nil
}

func (d *channelDataTxnHelper[TMemory]) CommitPreTick(p playback.Playback, cs *ChannelState[TMemory], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory]) error {
	// pre-effect
	if err := d.ProcessVolOps(p, cs); err != nil {
		return err
	}
	if err := d.ProcessNoteOps(p, cs); err != nil {
		return err
	}

	return nil
}

func (d *channelDataTxnHelper[TMemory]) CommitTick(p playback.Playback, cs *ChannelState[TMemory], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory]) error {
	if err := playback.DoEffect[TMemory](cs.ActiveEffect, cs, p, currentTick, lastTick); err != nil {
		return err
	}

	return nil
}

func (d *channelDataTxnHelper[TMemory]) CommitPostTick(p playback.Playback, cs *ChannelState[TMemory], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory]) error {
	// post-effect
	if err := d.ProcessVolOps(p, cs); err != nil {
		return err
	}
	if err := d.ProcessNoteOps(p, cs); err != nil {
		return err
	}

	return nil
}

func (d *channelDataTxnHelper[TMemory]) AddVolOp(op VolOp[TMemory]) {
	d.VolOps = append(d.VolOps, op)
}

func (d *channelDataTxnHelper[TMemory]) ProcessVolOps(p playback.Playback, cs *ChannelState[TMemory]) error {
	for _, op := range d.VolOps {
		if op == nil {
			continue
		}
		if err := op.Process(p, cs); err != nil {
			return err
		}
	}
	d.VolOps = nil

	return nil
}

func (d *channelDataTxnHelper[TMemory]) AddNoteOp(op NoteOp[TMemory]) {
	d.NoteOps = append(d.NoteOps, op)
}

func (d *channelDataTxnHelper[TMemory]) ProcessNoteOps(p playback.Playback, cs *ChannelState[TMemory]) error {
	for _, op := range d.NoteOps {
		if op == nil {
			continue
		}
		if err := op.Process(p, cs); err != nil {
			return err
		}
	}
	d.NoteOps = nil

	return nil
}

func (d *channelDataTxnHelper[TMemory]) ProcessData(out *ChannelDataActions, data song.ChannelData, s song.Data, cs *ChannelState[TMemory]) error {
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
				out.TargetNewNoteAction.Set(inst.GetNewNoteAction())
			}
		}
	}

	if note.IsInvalid(n) {
		out.TargetPeriod.Set(nil)
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
