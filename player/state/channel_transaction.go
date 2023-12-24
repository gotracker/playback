package state

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/op"
	"github.com/gotracker/playback/song"
)

type ChannelDataTransaction[TPeriod period.Period, TMemory any, TChannelData song.ChannelData] interface {
	GetChannelData() TChannelData
	SetData(data TChannelData, s song.Data, cs *ChannelState[TPeriod, TMemory, TChannelData]) error

	CommitPreRow(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData], semitoneSetterFactory SemitoneSetterFactory[TPeriod, TMemory, TChannelData]) error
	CommitRow(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData], semitoneSetterFactory SemitoneSetterFactory[TPeriod, TMemory, TChannelData]) error
	CommitPostRow(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData], semitoneSetterFactory SemitoneSetterFactory[TPeriod, TMemory, TChannelData]) error

	CommitPreTick(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TPeriod, TMemory, TChannelData]) error
	CommitTick(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TPeriod, TMemory, TChannelData]) error
	CommitPostTick(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TPeriod, TMemory, TChannelData]) error

	AddVolOp(op VolOp[TPeriod, TMemory, TChannelData])
	AddNoteOp(op NoteOp[TPeriod, TMemory, TChannelData])
}

type ChannelDataConverter[TPeriod period.Period, TMemory any, TChannelData song.ChannelData] func(out *op.ChannelTargets[TPeriod], data TChannelData, s song.Data, cs *ChannelState[TPeriod, TMemory, TChannelData]) error

type ChannelDataTxnHelper[TPeriod period.Period, TMemory any, TChannelData song.ChannelData] struct {
	Data     TChannelData
	Targeter playback.ChannelTargeter[TPeriod, TMemory, TChannelData]

	op.ChannelTargets[TPeriod]

	VolOps  []VolOp[TPeriod, TMemory, TChannelData]
	NoteOps []NoteOp[TPeriod, TMemory, TChannelData]
}

func NewChannelDataTxn[TPeriod period.Period, TMemory any, TChannelData song.ChannelData](targeter playback.ChannelTargeter[TPeriod, TMemory, TChannelData]) ChannelDataTransaction[TPeriod, TMemory, TChannelData] {
	return &ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]{
		Targeter: targeter,
	}
}

func (d *ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]) GetChannelData() TChannelData {
	return d.Data
}

func (d *ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]) SetData(cd TChannelData, s song.Data, cs *ChannelState[TPeriod, TMemory, TChannelData]) error {
	d.Data = cd
	return d.Targeter(&d.ChannelTargets, cd, s, cs)
}

func (d *ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]) CommitPreRow(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData], semitoneSetterFactory SemitoneSetterFactory[TPeriod, TMemory, TChannelData]) error {
	effects := playback.GetEffects[TPeriod, TMemory, TChannelData](cs.GetMemory(), d.Data)
	cs.SetActiveEffects(effects)
	onEff := p.GetOnEffect()
	for _, e := range effects {
		if onEff != nil {
			onEff(e)
		}
		if err := playback.EffectPreStart[TPeriod, TMemory](e, cs, p); err != nil {
			return err
		}
	}

	return nil
}

func (d *ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]) CommitRow(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData], semitoneSetterFactory SemitoneSetterFactory[TPeriod, TMemory, TChannelData]) error {
	target := cs.GetTargetState()
	if pos, ok := d.TargetPos.Get(); ok {
		target.Pos = pos
	}

	if inst, ok := d.TargetInst.Get(); ok {
		target.Instrument = inst
	}

	if period, ok := d.TargetPeriod.Get(); ok {
		target.Period = period
		cs.SetPortaTargetPeriod(period)
	}

	if st, ok := d.TargetStoredSemitone.Get(); ok {
		cs.SetStoredSemitone(st)
	}

	if nna, ok := d.TargetNewNoteAction.Get(); ok {
		cs.SetNewNoteAction(nna)
	}

	if v, ok := d.TargetVolume.Get(); ok {
		cs.GetActiveState().SetVolume(v)
		target.SetVolume(v)
	}

	na, targetTick := d.NoteAction.Get()
	cs.UseTargetPeriod = targetTick
	cs.SetNotePlayTick(targetTick, na, 0)

	if st, ok := d.NoteCalcST.Get(); ok {
		d.AddNoteOp(semitoneSetterFactory(st, func(p TPeriod) {
			target.Period = p
		}))
	}

	return nil
}

func (d *ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]) CommitPostRow(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData], semitoneSetterFactory SemitoneSetterFactory[TPeriod, TMemory, TChannelData]) error {
	return nil
}

func (d *ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]) CommitPreTick(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TPeriod, TMemory, TChannelData]) error {
	// pre-effect
	if err := d.ProcessVolOps(p, cs); err != nil {
		return err
	}
	if err := d.ProcessNoteOps(p, cs); err != nil {
		return err
	}

	return nil
}

func (d *ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]) CommitTick(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TPeriod, TMemory, TChannelData]) error {
	for _, e := range cs.GetActiveEffects() {
		if err := playback.DoEffect[TPeriod, TMemory](e, cs, p, currentTick, lastTick); err != nil {
			return err
		}
	}

	return nil
}

func (d *ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]) CommitPostTick(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData], currentTick int, lastTick bool, semitoneSetterFactory SemitoneSetterFactory[TPeriod, TMemory, TChannelData]) error {
	// post-effect
	if err := d.ProcessVolOps(p, cs); err != nil {
		return err
	}
	if err := d.ProcessNoteOps(p, cs); err != nil {
		return err
	}

	return nil
}

func (d *ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]) AddVolOp(op VolOp[TPeriod, TMemory, TChannelData]) {
	d.VolOps = append(d.VolOps, op)
}

func (d *ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]) ProcessVolOps(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData]) error {
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

func (d *ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]) AddNoteOp(op NoteOp[TPeriod, TMemory, TChannelData]) {
	d.NoteOps = append(d.NoteOps, op)
}

func (d *ChannelDataTxnHelper[TPeriod, TMemory, TChannelData]) ProcessNoteOps(p playback.Playback, cs *ChannelState[TPeriod, TMemory, TChannelData]) error {
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
