package state

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/heucuva/optional"
)

type ChannelDataTransaction[TChannelData any, TChannelState playback.ChannelState] interface {
	GetData() *TChannelData
	SetData(data *TChannelData, cs *TChannelState) error

	CommitPreRow(p playback.Playback, cs *TChannelState) error
	CommitRow(p playback.Playback, cs *TChannelState) error
	CommitPostRow(p playback.Playback, cs *TChannelState) error

	CommitPreTick(p playback.Playback, cs *TChannelState, currentTick int, lastTick bool) error
	CommitTick(p playback.Playback, cs *TChannelState, currentTick int, lastTick bool) error
	CommitPostTick(p playback.Playback, cs *TChannelState, currentTick int, lastTick bool) error

	AddVolOp(op VolOp[TChannelState])
	AddNoteOp(op NoteOp[TChannelState])
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

type ChannelDataConverter[TChannelData any, TChannelState playback.ChannelState] interface {
	Process(out *ChannelDataActions, data *TChannelData, cs *TChannelState) error
}

type ChannelDataTxnHelper[TChannelData any, TChannelState playback.ChannelState, TChannelDataConverter ChannelDataConverter[TChannelData, TChannelState]] struct {
	Data *TChannelData

	ChannelDataActions

	VolOps  []VolOp[TChannelState]
	NoteOps []NoteOp[TChannelState]
}

func (d *ChannelDataTxnHelper[TChannelData, TChannelState, TChannelDataConverter]) GetData() *TChannelData {
	return d.Data
}

func (d *ChannelDataTxnHelper[TChannelData, TChannelState, TChannelDataConverter]) SetData(cd *TChannelData, cs *TChannelState) error {
	d.Data = cd

	var converter TChannelDataConverter
	return converter.Process(&d.ChannelDataActions, cd, cs)
}

func (d *ChannelDataTxnHelper[TChannelData, TChannelState, TChannelDataConverter]) CommitPreRow(p playback.Playback, cs *TChannelState) error {
	return nil
}

func (d *ChannelDataTxnHelper[TChannelData, TChannelState, TChannelDataConverter]) CommitRow(p playback.Playback, cs *TChannelState) error {
	return nil
}

func (d *ChannelDataTxnHelper[TChannelData, TChannelState, TChannelDataConverter]) CommitPostRow(p playback.Playback, cs *TChannelState) error {
	return nil
}

func (d *ChannelDataTxnHelper[TChannelData, TChannelState, TChannelDataConverter]) CommitPreTick(p playback.Playback, cs *TChannelState, currentTick int, lastTick bool) error {
	// pre-effect
	if err := d.ProcessVolOps(p, cs); err != nil {
		return err
	}
	if err := d.ProcessNoteOps(p, cs); err != nil {
		return err
	}

	return nil
}

func (d *ChannelDataTxnHelper[TChannelData, TChannelState, TChannelDataConverter]) CommitTick(p playback.Playback, cs *TChannelState, currentTick int, lastTick bool) error {
	pce := any(cs).(playback.ChannelEffect[TChannelState])
	e := pce.GetActiveEffect()
	if err := playback.DoEffect(e, cs, p, currentTick, lastTick); err != nil {
		return err
	}

	return nil
}

func (d *ChannelDataTxnHelper[TChannelData, TChannelState, TChannelDataConverter]) CommitPostTick(p playback.Playback, cs *TChannelState, currentTick int, lastTick bool) error {
	// post-effect
	if err := d.ProcessVolOps(p, cs); err != nil {
		return err
	}
	if err := d.ProcessNoteOps(p, cs); err != nil {
		return err
	}

	return nil
}

func (d *ChannelDataTxnHelper[TChannelData, TChannelState, TChannelDataConverter]) AddVolOp(op VolOp[TChannelState]) {
	d.VolOps = append(d.VolOps, op)
}

func (d *ChannelDataTxnHelper[TChannelData, TChannelState, TChannelDataConverter]) ProcessVolOps(p playback.Playback, cs *TChannelState) error {
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

func (d *ChannelDataTxnHelper[TChannelData, TChannelState, TChannelDataConverter]) AddNoteOp(op NoteOp[TChannelState]) {
	d.NoteOps = append(d.NoteOps, op)
}

func (d *ChannelDataTxnHelper[TChannelData, TChannelState, TChannelDataConverter]) ProcessNoteOps(p playback.Playback, cs *TChannelState) error {
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
