package channel

import (
	"fmt"
	"strings"

	xmfile "github.com/gotracker/goaudiofile/music/tracked/xm"

	"github.com/gotracker/playback"
	xmNote "github.com/gotracker/playback/format/xm/note"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/player/machine/instruction"
	"github.com/gotracker/playback/song"
)

type Command uint8

func (c Command) ToRune() rune {
	switch {
	case c <= 9:
		return '0' + rune(c)
	case c >= 10 && c < 36:
		return 'A' + rune(c-10)
	default:
		panic("effect out of range")
	}
}

// DataEffect is the type of a channel's EffectParameter value
type DataEffect uint8

// Data is the data for the channel
type Data[TPeriod period.Period] struct {
	What            xmfile.ChannelFlags
	Note            uint8
	Instrument      uint8
	Volume          xmVolume.VolEffect
	Effect          Command
	EffectParameter DataEffect
}

// HasNote returns true if there exists a note on the channel
func (d Data[TPeriod]) HasNote() bool {
	return d.What.HasNote()
}

// GetNote returns the note for the channel
func (d Data[TPeriod]) GetNote() note.Note {
	return xmNote.FromXmNote(d.Note)
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
	if !d.What.HasVolume() {
		return false
	}

	return d.Volume.IsVolume()
}

// GetVolume returns the volume for the channel
func (d Data[TPeriod]) GetVolumeGeneric() volume.Volume {
	return d.Volume.Volume()
}

func (d Data[TPeriod]) GetVolume() xmVolume.XmVolume {
	return xmVolume.XmVolume(d.Volume.Volume() * 64)
}

// HasCommand returns true if there exists a command on the channel
func (d Data[TPeriod]) HasCommand() bool {
	if d.What.HasEffect() || d.What.HasEffectParameter() {
		return true
	}

	if d.What.HasVolume() {
		return !d.Volume.IsVolume()
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
		return "== "
	case note.SpecialTypeNormal:
		return n.String()
	default:
		return "???"
	}
}

func (d Data[TPeriod]) String() string {
	pieces := []string{
		"...", // note
		"  ",  // inst
		"..",  // vol
		"...", // eff
	}

	if d.HasNote() {
		pieces[0] = d.getNoteString(d.GetNote())
	}
	if d.HasInstrument() {
		pieces[1] = fmt.Sprintf("%2X", d.Instrument)
	}
	if d.HasVolume() {
		pieces[2] = fmt.Sprintf("%02X", d.Volume)
	}
	if d.HasCommand() {
		pieces[3] = fmt.Sprintf("%c%02X", d.Effect.ToRune(), d.EffectParameter)
	}
	return strings.Join(pieces, " ")
}

func (d Data[TPeriod]) ShortString() string {
	if d.HasNote() {
		return d.getNoteString(d.GetNote())
	}
	return "..."
}

func (d Data[TPeriod]) ToInstructions(m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], ch index.Channel, songData song.Data) ([]instruction.Instruction, error) {
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
