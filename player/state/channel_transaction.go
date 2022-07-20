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

type ChannelDataTransaction[TMemory, TChannelData any] interface {
	GetData() *TChannelData
	SetData(data *TChannelData, s song.Data, cs *ChannelState[TMemory, TChannelData]) error

	CommitPreRow(p playback.Playback, cs *ChannelState[TMemory, TChannelData], semitoneSetterFactory SemitoneSetterFactory[TMemory, TChannelData]) error
	CommitRow(p playback.Playback, cs *ChannelState[TMemory, TChannelData], semitoneSetterFactory SemitoneSetterFactory[TMemory, TChannelData]) error
	CommitPostRow(p playback.Playback, cs *ChannelState[TMemory, TChannelData], semitoneSetterFactory SemitoneSetterFactory[TMemory, TChannelData]) error

	CommitPreTick(p playback.Playback, cs *ChannelState[TMemory, TChannelData], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory, TChannelData]) error
	CommitTick(p playback.Playback, cs *ChannelState[TMemory, TChannelData], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory, TChannelData]) error
	CommitPostTick(p playback.Playback, cs *ChannelState[TMemory, TChannelData], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory, TChannelData]) error

	AddVolOp(op VolOp[TMemory, TChannelData])
	AddNoteOp(op NoteOp[TMemory, TChannelData])
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

type ChannelDataConverter[TMemory, TChannelData any] interface {
	Process(out *ChannelDataActions, data *TChannelData, s song.Data, cs *ChannelState[TMemory, TChannelData]) error
}

type ChannelDataTxnHelper[TMemory, TChannelData any, TChannelDataConverter ChannelDataConverter[TMemory, TChannelData]] struct {
	Data *TChannelData

	ChannelDataActions

	VolOps  []VolOp[TMemory, TChannelData]
	NoteOps []NoteOp[TMemory, TChannelData]
}

func (d *ChannelDataTxnHelper[TMemory, TChannelData, TChannelDataConverter]) GetData() *TChannelData {
	return d.Data
}

func (d *ChannelDataTxnHelper[TMemory, TChannelData, TChannelDataConverter]) SetData(cd *TChannelData, s song.Data, cs *ChannelState[TMemory, TChannelData]) error {
	d.Data = cd

	var converter TChannelDataConverter
	return converter.Process(&d.ChannelDataActions, cd, s, cs)
}

func (d *ChannelDataTxnHelper[TMemory, TChannelData, TChannelDataConverter]) CommitPreRow(p playback.Playback, cs *ChannelState[TMemory, TChannelData], semitoneSetterFactory SemitoneSetterFactory[TMemory, TChannelData]) error {
	return nil
}

func (d *ChannelDataTxnHelper[TMemory, TChannelData, TChannelDataConverter]) CommitRow(p playback.Playback, cs *ChannelState[TMemory, TChannelData], semitoneSetterFactory SemitoneSetterFactory[TMemory, TChannelData]) error {
	return nil
}

func (d *ChannelDataTxnHelper[TMemory, TChannelData, TChannelDataConverter]) CommitPostRow(p playback.Playback, cs *ChannelState[TMemory, TChannelData], semitoneSetterFactory SemitoneSetterFactory[TMemory, TChannelData]) error {
	return nil
}

func (d *ChannelDataTxnHelper[TMemory, TChannelData, TChannelDataConverter]) CommitPreTick(p playback.Playback, cs *ChannelState[TMemory, TChannelData], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory, TChannelData]) error {
	// pre-effect
	if err := d.ProcessVolOps(p, cs); err != nil {
		return err
	}
	if err := d.ProcessNoteOps(p, cs); err != nil {
		return err
	}

	return nil
}

func (d *ChannelDataTxnHelper[TMemory, TChannelData, TChannelDataConverter]) CommitTick(p playback.Playback, cs *ChannelState[TMemory, TChannelData], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory, TChannelData]) error {
	if err := playback.DoEffect[TMemory, TChannelData](cs.ActiveEffect, cs, p, currentTick, lastTick); err != nil {
		return err
	}

	return nil
}

func (d *ChannelDataTxnHelper[TMemory, TChannelData, TChannelDataConverter]) CommitPostTick(p playback.Playback, cs *ChannelState[TMemory, TChannelData], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TMemory, TChannelData]) error {
	// post-effect
	if err := d.ProcessVolOps(p, cs); err != nil {
		return err
	}
	if err := d.ProcessNoteOps(p, cs); err != nil {
		return err
	}

	return nil
}

func (d *ChannelDataTxnHelper[TMemory, TChannelData, TChannelDataConverter]) AddVolOp(op VolOp[TMemory, TChannelData]) {
	d.VolOps = append(d.VolOps, op)
}

func (d *ChannelDataTxnHelper[TMemory, TChannelData, TChannelDataConverter]) ProcessVolOps(p playback.Playback, cs *ChannelState[TMemory, TChannelData]) error {
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

func (d *ChannelDataTxnHelper[TMemory, TChannelData, TChannelDataConverter]) AddNoteOp(op NoteOp[TMemory, TChannelData]) {
	d.NoteOps = append(d.NoteOps, op)
}

func (d *ChannelDataTxnHelper[TMemory, TChannelData, TChannelDataConverter]) ProcessNoteOps(p playback.Playback, cs *ChannelState[TMemory, TChannelData]) error {
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
