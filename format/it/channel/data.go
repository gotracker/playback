package channel

import (
	"fmt"
	"strings"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"

	"github.com/gotracker/playback"
	itNote "github.com/gotracker/playback/format/it/note"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/player/machine/instruction"
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
func (d Data[TPeriod]) GetInstrument() int {
	return int(d.Instrument)
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
	case note.SpecialTypeFadeout:
		return "vvv"
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
