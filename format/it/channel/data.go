package channel

import (
	"fmt"
	"strings"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback"
	itNote "github.com/gotracker/playback/format/it/note"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/op"
	"github.com/gotracker/playback/song"
)

const MaxTotalChannels = 64

type Command uint8

func (c Command) ToRune() rune {
	switch {
	case c > 0 && c <= 26:
		return '@' + rune(c)
	default:
		panic("effect out of range")
	}
}

// DataEffect is the type of a channel's EffectParameter value
type DataEffect uint8

// Data is the data for the channel
type Data struct {
	What            itfile.ChannelDataFlags
	Note            itfile.Note
	Instrument      uint8
	VolPan          uint8
	Effect          Command
	EffectParameter DataEffect
}

// HasNote returns true if there exists a note on the channel
func (d Data) HasNote() bool {
	return d.What.HasNote()
}

// GetNote returns the note for the channel
func (d Data) GetNote() note.Note {
	return itNote.FromItNote(d.Note)
}

// HasInstrument returns true if there exists an instrument on the channel
func (d Data) HasInstrument() bool {
	return d.What.HasInstrument()
}

// GetInstrument returns the instrument for the channel
func (d Data) GetInstrument(stmem note.Semitone) instrument.ID {
	st := stmem
	if d.HasNote() {
		n := d.GetNote()
		if nn, ok := n.(note.Normal); ok {
			st = note.Semitone(nn)
		}
	}
	return SampleID{
		InstID:   d.Instrument,
		Semitone: st,
	}
}

// HasVolume returns true if there exists a volume on the channel
func (d Data) HasVolume() bool {
	if !d.What.HasVolPan() {
		return false
	}

	v := d.VolPan
	return v <= 64
}

// GetVolume returns the volume for the channel
func (d Data) GetVolume() volume.Volume {
	return itVolume.FromVolPan(d.VolPan)
}

// HasCommand returns true if there exists a effect on the channel
func (d Data) HasCommand() bool {
	if d.What.HasCommand() {
		return true
	}

	if d.What.HasVolPan() {
		return d.VolPan > 64
	}

	return false
}

// Channel returns the channel ID for the channel
func (d Data) Channel() uint8 {
	return 0
}

func (d Data) GetEffects(mem *Memory, periodType period.Period) []playback.Effect {
	switch periodType.(type) {
	case period.Linear:
		if e := EffectFactory[period.Linear](mem, d); e != nil {
			return []playback.Effect{e}
		}
	case period.Amiga:
		if e := EffectFactory[period.Amiga](mem, d); e != nil {
			return []playback.Effect{e}
		}
	default:
		panic("unhandled period type")
	}
	return nil
}

func (Data) getNoteString(n note.Note) string {
	switch note.Type(n) {
	case note.SpecialTypeRelease:
		return "==="
	case note.SpecialTypeStop:
		return "^^^"
	case note.SpecialTypeNormal:
		return n.String()
	default:
		return "???"
	}
}

func (d Data) String() string {
	pieces := []string{
		"...", // note
		"..",  // inst
		"..",  // vol
		"...", // eff
	}
	if d.HasNote() {
		pieces[0] = d.getNoteString(d.GetNote())
	}
	if d.HasInstrument() {
		pieces[1] = fmt.Sprintf("%02X", d.Instrument)
	}
	if d.HasVolume() {
		pieces[2] = fmt.Sprintf("%02X", d.VolPan)
	}
	if d.HasCommand() && d.Effect != 0 {
		pieces[3] = fmt.Sprintf("%c%02X", d.Effect.ToRune(), d.EffectParameter)
	}
	return strings.Join(pieces, " ")
}

func (d Data) ShortString() string {
	if d.HasNote() {
		return d.GetNote().String()
	}
	return "..."
}

func GetTargetsFromData[TPeriod period.Period](out *op.ChannelTargets[TPeriod], d Data, s song.Data, cs playback.Channel[TPeriod, Memory, Data]) error {
	var n note.Note = note.EmptyNote{}
	inst := cs.GetActiveState().Instrument
	prevInst := inst

	if d.HasNote() || d.HasInstrument() {
		instID := d.GetInstrument(cs.GetNoteSemitone())
		n = d.GetNote()
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
		out.TargetPeriod.Reset()
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

	if d.HasVolume() {
		v := d.GetVolume()
		if v == volume.VolumeUseInstVol {
			if inst != nil {
				v = inst.GetDefaultVolume()
			}
		}
		out.TargetVolume.Set(v)
	}

	return nil
}
