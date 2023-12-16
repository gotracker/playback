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

type ChannelDataConverter[TMemory any] interface {
	Process(out *ChannelDataActions, data song.ChannelData, s song.Data, cs *ChannelState[TMemory]) error
}

type ChannelDataTxnHelper[TMemory any, TChannelDataConverter ChannelDataConverter[TMemory]] struct {
	Data song.ChannelData

	ChannelDataActions

	VolOps  []VolOp[TMemory]
	NoteOps []NoteOp[TMemory]
}

func (d *ChannelDataTxnHelper[TMemory, TChannelDataConverter]) GetData() song.ChannelData {
	return d.Data
}

func (d *ChannelDataTxnHelper[TMemory, TChannelDataConverter]) SetData(cd song.ChannelData, s song.Data, cs *ChannelState[TMemory]) error {
	d.Data = cd

	var converter TChannelDataConverter
	return converter.Process(&d.ChannelDataActions, cd, s, cs)
}

func (d *ChannelDataTxnHelper[TMemory, TChannelDataConverter]) CommitPreRow(p playback.Playback, cs *ChannelState[TMemory], semitoneSetterFactory SemitoneSetterFactory[TMemory]) error {
	return nil
}

func (d *ChannelDataTxnHelper[TMemory, TChannelDataConverter]) CommitRow(p playback.Playback, cs *ChannelState[TMemory], semitoneSetterFactory SemitoneSetterFactory[TMemory]) error {
	return nil
}

func (d *ChannelDataTxnHelper[TMemory, TChannelDataConverter]) CommitPostRow(p playback.Playback, cs *ChannelState[TMemory], semitoneSetterFactory SemitoneSetterFactory[TMemory]) error {
	return nil
}

func (d *ChannelDataTxnHelper[TMemory, TChannelDataConverter]) CommitPreTick(p playback.Playback, cs *ChannelState[TMemory], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory]) error {
	// pre-effect
	if err := d.ProcessVolOps(p, cs); err != nil {
		return err
	}
	if err := d.ProcessNoteOps(p, cs); err != nil {
		return err
	}

	return nil
}

func (d *ChannelDataTxnHelper[TMemory, TChannelDataConverter]) CommitTick(p playback.Playback, cs *ChannelState[TMemory], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory]) error {
	if err := playback.DoEffect[TMemory](cs.ActiveEffect, cs, p, currentTick, lastTick); err != nil {
		return err
	}

	return nil
}

func (d *ChannelDataTxnHelper[TMemory, TChannelDataConverter]) CommitPostTick(p playback.Playback, cs *ChannelState[TMemory], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory]) error {
	// post-effect
	if err := d.ProcessVolOps(p, cs); err != nil {
		return err
	}
	if err := d.ProcessNoteOps(p, cs); err != nil {
		return err
	}

	return nil
}

func (d *ChannelDataTxnHelper[TMemory, TChannelDataConverter]) AddVolOp(op VolOp[TMemory]) {
	d.VolOps = append(d.VolOps, op)
}

func (d *ChannelDataTxnHelper[TMemory, TChannelDataConverter]) ProcessVolOps(p playback.Playback, cs *ChannelState[TMemory]) error {
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

func (d *ChannelDataTxnHelper[TMemory, TChannelDataConverter]) AddNoteOp(op NoteOp[TMemory]) {
	d.NoteOps = append(d.NoteOps, op)
}

func (d *ChannelDataTxnHelper[TMemory, TChannelDataConverter]) ProcessNoteOps(p playback.Playback, cs *ChannelState[TMemory]) error {
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
