package channel

import (
	"fmt"
	"strings"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback"
	itNote "github.com/gotracker/playback/format/it/note"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/player/machine/instruction"
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
type Data[TPeriod period.Period] struct {
	What            itfile.ChannelDataFlags
	Note            itfile.Note
	Instrument      uint8
	VolPan          uint8
	Effect          Command
	EffectParameter DataEffect
}

// HasNote returns true if there exists a note on the channel
func (d Data[TPeriod]) HasNote() bool {
	return d.What.HasNote()
}

// GetNote returns the note for the channel
func (d Data[TPeriod]) GetNote() note.Note {
	return itNote.FromItNote(d.Note)
}

// HasInstrument returns true if there exists an instrument on the channel
func (d Data[TPeriod]) HasInstrument() bool {
	return d.What.HasInstrument()
}

// GetInstrument returns the instrument for the channel
func (d Data[TPeriod]) GetInstrument(stmem note.Semitone) instrument.ID {
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
func (d Data[TPeriod]) HasVolume() bool {
	if !d.What.HasVolPan() {
		return false
	}

	v := d.VolPan
	return v <= 64
}

// GetVolume returns the volume for the channel
func (d Data[TPeriod]) GetVolumeGeneric() volume.Volume {
	return itVolume.FromVolPan(d.VolPan)
}

func (d Data[TPeriod]) GetVolume() itVolume.Volume {
	return itVolume.Volume(d.VolPan)
}

// HasCommand returns true if there exists a effect on the channel
func (d Data[TPeriod]) HasCommand() bool {
	if d.What.HasCommand() {
		return true
	}

	if d.What.HasVolPan() {
		return d.VolPan > 64
	}

	return false
}

// Channel returns the channel ID for the channel
func (d Data[TPeriod]) Channel() uint8 {
	return 0
}

func (d Data[TPeriod]) GetEffects(mem *Memory) []playback.Effect {
	if e := EffectFactory[TPeriod](mem, d); e != nil {
		return []playback.Effect{e}
	}
	return nil
}

func (Data[TPeriod]) getNoteString(n note.Note) string {
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

func (d Data[TPeriod]) String() string {
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

func (d Data[TPeriod]) ShortString() string {
	if d.HasNote() {
		return d.GetNote().String()
	}
	return "..."
}

func (d Data[TPeriod]) ToInstructions(m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], ch index.Channel, songData song.Data) ([]instruction.Instruction, error) {
	var instructions []instruction.Instruction

	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return nil, err
	}

	if e := EffectFactory[TPeriod](mem, d); e != nil {
		instructions = append(instructions, e)
	}

	return instructions, nil
}

func GetTargetsFromData[TPeriod period.Period](out *op.ChannelTargets[TPeriod, itVolume.Volume, itPanning.Panning], d Data[TPeriod], s song.Data, cs playback.Channel[TPeriod, *Memory, Data[TPeriod], itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
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
					out.TargetVolume.Set(itVolume.ToItVolume(inst.GetDefaultVolumeGeneric()))
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
		if v.IsUseInstrumentVol() {
			if inst != nil {
				v = itVolume.ToItVolume(inst.GetDefaultVolumeGeneric())
			}
		}
		out.TargetVolume.Set(v)
	}

	return nil
}
