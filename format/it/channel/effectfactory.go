package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
)

type EffectIT = playback.Effect

// VolEff is a combined effect that includes a volume effect and a standard effect
type VolEff[TPeriod period.Period] struct {
	playback.CombinedEffect[TPeriod, Memory]
	eff EffectIT
}

func (e VolEff[TPeriod]) String() string {
	if e.eff == nil {
		return "..."
	}
	return fmt.Sprint(e.eff)
}

// Factory produces an effect for the provided channel pattern data
func EffectFactory[TPeriod period.Period](mem *Memory, data song.ChannelData) EffectIT {
	if data == nil {
		return nil
	}

	d, _ := data.(Data)

	if !d.What.HasCommand() && !d.What.HasVolPan() {
		return nil
	}

	var eff VolEff[TPeriod]
	if d.What.HasVolPan() {
		ve := volPanEffectFactory[TPeriod](mem, d.VolPan)
		if ve != nil {
			eff.Effects = append(eff.Effects, ve)
		}
	}

	if e := standardEffectFactory[TPeriod](mem, d); e != nil {
		eff.Effects = append(eff.Effects, e)
		eff.eff = e
	}

	switch len(eff.Effects) {
	case 0:
		return nil
	case 1:
		return eff.Effects[0]
	default:
		return &eff
	}
}

func standardEffectFactory[TPeriod period.Period](mem *Memory, data Data) EffectIT {
	switch data.Effect + '@' {
	case '@': // unused
		return nil
	case 'A': // Set Speed
		return SetSpeed[TPeriod](data.EffectParameter)
	case 'B': // Pattern Jump
		return OrderJump[TPeriod](data.EffectParameter)
	case 'C': // Pattern Break
		return RowJump[TPeriod](data.EffectParameter)
	case 'D': // Volume Slide / Fine Volume Slide
		return volumeSlideFactory[TPeriod](mem, data.Effect, data.EffectParameter)
	case 'E': // Porta Down/Fine Porta Down/Xtra Fine Porta
		xx := mem.PortaDown(DataEffect(data.EffectParameter))
		x := xx >> 4
		if x == 0x0F {
			return FinePortaDown[TPeriod](xx)
		} else if x == 0x0E {
			return ExtraFinePortaDown[TPeriod](xx)
		}
		return PortaDown[TPeriod](data.EffectParameter)
	case 'F': // Porta Up/Fine Porta Up/Extra Fine Porta Down
		xx := mem.PortaUp(DataEffect(data.EffectParameter))
		x := xx >> 4
		if x == 0x0F {
			return FinePortaUp[TPeriod](xx)
		} else if x == 0x0E {
			return ExtraFinePortaUp[TPeriod](xx)
		}
		return PortaUp[TPeriod](data.EffectParameter)
	case 'G': // Porta to note
		return PortaToNote[TPeriod](data.EffectParameter)
	case 'H': // Vibrato
		return Vibrato[TPeriod](data.EffectParameter)
	case 'I': // Tremor
		return Tremor[TPeriod](data.EffectParameter)
	case 'J': // Arpeggio
		return Arpeggio[TPeriod](data.EffectParameter)
	case 'K': // Vibrato+Volume Slide
		return NewVibratoVolumeSlide[TPeriod](mem, data.Effect, data.EffectParameter)
	case 'L': // Porta+Volume Slide
		return NewPortaVolumeSlide[TPeriod](mem, data.Effect, data.EffectParameter)
	case 'M': // Set Channel Volume
		return SetChannelVolume[TPeriod](data.EffectParameter)
	case 'N': // Channel Volume Slide
		return ChannelVolumeSlide[TPeriod](data.EffectParameter)
	case 'O': // Sample Offset
		return SampleOffset[TPeriod](data.EffectParameter)
	case 'P': // Panning Slide
		//return panningSlideFactory(mem, data.Effect, data.EffectParameter)
	case 'Q': // Retrig + Volume Slide
		return RetrigVolumeSlide[TPeriod](data.EffectParameter)
	case 'R': // Tremolo
		return Tremolo[TPeriod](data.EffectParameter)
	case 'S': // Special
		return specialEffect[TPeriod](data)
	case 'T': // Set Tempo
		return SetTempo[TPeriod](data.EffectParameter)
	case 'U': // Fine Vibrato
		return FineVibrato[TPeriod](data.EffectParameter)
	case 'V': // Global Volume
		return SetGlobalVolume[TPeriod](data.EffectParameter)
	case 'W': // Global Volume Slide
		return GlobalVolumeSlide[TPeriod](data.EffectParameter)
	case 'X': // Set Pan Position
		return SetPanPosition[TPeriod](data.EffectParameter)
	case 'Y': // Panbrello
		//return Panbrello[TPeriod](data.EffectParameter)
	case 'Z': // MIDI Macro
		return nil // TODO: MIDIMacro
	default:
	}
	return UnhandledCommand[TPeriod]{Command: data.Effect, Info: data.EffectParameter}
}

func specialEffect[TPeriod period.Period](data Data) EffectIT {
	switch data.EffectParameter >> 4 {
	case 0x0: // unused
		return nil
	//case 0x1: // Set Glissando on/off

	case 0x2: // Set FineTune
		return SetFinetune[TPeriod](data.EffectParameter)
	case 0x3: // Set Vibrato Waveform
		return SetVibratoWaveform[TPeriod](data.EffectParameter)
	case 0x4: // Set Tremolo Waveform
		return SetTremoloWaveform[TPeriod](data.EffectParameter)
	case 0x5: // Set Panbrello Waveform
		return SetPanbrelloWaveform[TPeriod](data.EffectParameter)
	case 0x6: // Fine Pattern Delay
		return FinePatternDelay[TPeriod](data.EffectParameter)
	case 0x7: // special note operations
		return specialNoteEffects[TPeriod](data)
	case 0x8: // Set Coarse Pan Position
		return SetCoarsePanPosition[TPeriod](data.EffectParameter)
	case 0x9: // Sound Control
		return soundControlEffect[TPeriod](data)
	case 0xA: // High Offset
		return HighOffset[TPeriod](data.EffectParameter)
	case 0xB: // Pattern Loop
		return PatternLoop[TPeriod](data.EffectParameter)
	case 0xC: // Note Cut
		return NoteCut[TPeriod](data.EffectParameter)
	case 0xD: // Note Delay
		return NoteDelay[TPeriod](data.EffectParameter)
	case 0xE: // Pattern Delay
		return PatternDelay[TPeriod](data.EffectParameter)
	case 0xF: // Set Active Macro
		return nil // TODO: SetActiveMacro
	default:
	}
	return UnhandledCommand[TPeriod]{Command: data.Effect, Info: data.EffectParameter}
}

func specialNoteEffects[TPeriod period.Period](data Data) EffectIT {
	switch data.EffectParameter & 0xf {
	case 0x0: // Past Note Cut
		return PastNoteCut[TPeriod](data.EffectParameter)
	case 0x1: // Past Note Off
		return PastNoteOff[TPeriod](data.EffectParameter)
	case 0x2: // Past Note Fade
		return PastNoteFade[TPeriod](data.EffectParameter)
	case 0x3: // New Note Action: Note Cut
		return NewNoteActionNoteCut[TPeriod](data.EffectParameter)
	case 0x4: // New Note Action: Note Continue
		return NewNoteActionNoteContinue[TPeriod](data.EffectParameter)
	case 0x5: // New Note Action: Note Off
		return NewNoteActionNoteOff[TPeriod](data.EffectParameter)
	case 0x6: // New Note Action: Note Fade
		return NewNoteActionNoteFade[TPeriod](data.EffectParameter)
	case 0x7: // Volume Envelope Off
		return VolumeEnvelopeOff[TPeriod](data.EffectParameter)
	case 0x8: // Volume Envelope On
		return VolumeEnvelopeOn[TPeriod](data.EffectParameter)
	case 0x9: // Panning Envelope Off
		return PanningEnvelopeOff[TPeriod](data.EffectParameter)
	case 0xA: // Panning Envelope On
		return PanningEnvelopeOn[TPeriod](data.EffectParameter)
	case 0xB: // Pitch Envelope Off
		return PitchEnvelopeOff[TPeriod](data.EffectParameter)
	case 0xC: // Pitch Envelope On
		return PitchEnvelopeOn[TPeriod](data.EffectParameter)
	case 0xD, 0xE, 0xF: // unused
		return nil
	}
	return UnhandledCommand[TPeriod]{Command: data.Effect, Info: data.EffectParameter}
}

func volumeSlideFactory[TPeriod period.Period](mem *Memory, cd Command, ce DataEffect) EffectIT {
	x, y := mem.VolumeSlide(DataEffect(ce))
	switch {
	case x == 0:
		return VolumeSlideDown[TPeriod](ce)
	case y == 0:
		return VolumeSlideUp[TPeriod](ce)
	case x == 0x0f:
		return FineVolumeSlideDown[TPeriod](ce)
	case y == 0x0f:
		return FineVolumeSlideUp[TPeriod](ce)
	}
	// There is a chance that a volume slide command is set with an invalid
	// value or is 00, in which case the memory might have the invalid value,
	// so we need to handle it by deferring to using a no-op instead of a
	// VolumeSlideDown
	return nil
}

func soundControlEffect[TPeriod period.Period](data Data) EffectIT {
	switch data.EffectParameter & 0xF {
	case 0x0: // Surround Off
	case 0x1: // Surround On
		// only S91 is supported directly by IT
		return nil // TODO: SurroundOn
	case 0x8: // Reverb Off
	case 0x9: // Reverb On
	case 0xA: // Center Surround
	case 0xB: // Quad Surround
	case 0xC: // Global Filters
	case 0xD: // Local Filters
	case 0xE: // Play Forward
	case 0xF: // Play Backward
	}
	return UnhandledCommand[TPeriod]{Command: data.Effect, Info: data.EffectParameter}
}
