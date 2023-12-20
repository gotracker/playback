package effect

import (
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/period"
)

func standardEffectFactory[TPeriod period.Period](mem *channel.Memory, cd channel.Data) EffectXM {
	if !cd.What.HasEffect() && !cd.What.HasEffectParameter() {
		return nil
	}

	switch cd.Effect {
	case 0x00: // Arpeggio
		return Arpeggio[TPeriod](cd.EffectParameter)
	case 0x01: // Porta up
		return PortaUp[TPeriod](cd.EffectParameter)
	case 0x02: // Porta down
		return PortaDown[TPeriod](cd.EffectParameter)
	case 0x03: // Tone porta
		return PortaToNote[TPeriod](cd.EffectParameter)
	case 0x04: // Vibrato
		return Vibrato[TPeriod](cd.EffectParameter)
	case 0x05: // Tone porta + Volume slide
		return NewPortaVolumeSlide[TPeriod](cd.EffectParameter)
	case 0x06: // Vibrato + Volume slide
		return NewVibratoVolumeSlide[TPeriod](cd.EffectParameter)
	case 0x07: // Tremolo
		return Tremolo[TPeriod](cd.EffectParameter)
	case 0x08: // Set (fine) panning
		return SetPanPosition[TPeriod](cd.EffectParameter)
	case 0x09: // Sample offset
		return SampleOffset[TPeriod](cd.EffectParameter)
	case 0x0A: // Volume slide
		return VolumeSlide[TPeriod](cd.EffectParameter)
	case 0x0B: // Position jump
		return OrderJump[TPeriod](cd.EffectParameter)
	case 0x0C: // Set volume
		return SetVolume[TPeriod](cd.EffectParameter)
	case 0x0D: // Pattern break
		return RowJump[TPeriod](cd.EffectParameter)
	case 0x0E: // extra...
		return specialEffectFactory[TPeriod](mem, cd.Effect, cd.EffectParameter)
	case 0x0F: // Set tempo/BPM
		if cd.EffectParameter < 0x20 {
			return SetSpeed[TPeriod](cd.EffectParameter)
		}
		return SetTempo[TPeriod](cd.EffectParameter)
	case 0x10: // Set global volume
		return SetGlobalVolume[TPeriod](cd.EffectParameter)
	case 0x11: // Global volume slide
		return GlobalVolumeSlide[TPeriod](cd.EffectParameter)

	case 0x15: // Set envelope position
		return SetEnvelopePosition[TPeriod](cd.EffectParameter)

	case 0x19: // Panning slide
		return PanSlide[TPeriod](cd.EffectParameter)

	case 0x1B: // Multi retrig note
		return RetrigVolumeSlide[TPeriod](cd.EffectParameter)

	case 0x1D: // Tremor
		return Tremor[TPeriod](cd.EffectParameter)

	case 0x21: // Extra fine porta commands
		return extraFinePortaEffectFactory[TPeriod](mem, cd.Effect, cd.EffectParameter)
	}
	return UnhandledCommand[TPeriod]{Command: cd.Effect, Info: cd.EffectParameter}
}

func extraFinePortaEffectFactory[TPeriod period.Period](mem *channel.Memory, ce channel.Command, cp channel.DataEffect) EffectXM {
	switch cp >> 4 {
	case 0x0: // none
		return nil
	case 0x1: // Extra fine porta up
		return ExtraFinePortaUp[TPeriod](cp)
	case 0x2: // Extra fine porta down
		return ExtraFinePortaDown[TPeriod](cp)
	}
	return UnhandledCommand[TPeriod]{Command: ce, Info: cp}
}

func specialEffectFactory[TPeriod period.Period](mem *channel.Memory, ce channel.Command, cp channel.DataEffect) EffectXM {
	switch cp >> 4 {
	case 0x1: // Fine porta up
		return FinePortaUp[TPeriod](cp)
	case 0x2: // Fine porta down
		return FinePortaDown[TPeriod](cp)
	//case 0x3: // Set glissando control

	case 0x4: // Set vibrato control
		return SetVibratoWaveform[TPeriod](cp)
	case 0x5: // Set finetune
		return SetFinetune[TPeriod](cp)
	case 0x6: // Set loop begin/loop
		return PatternLoop[TPeriod](cp)
	case 0x7: // Set tremolo control
		return SetTremoloWaveform[TPeriod](cp)
	case 0x8: // Set coarse panning
		return SetCoarsePanPosition[TPeriod](cp)
	case 0x9: // Retrig note
		return RetriggerNote[TPeriod](cp)
	case 0xA: // Fine volume slide up
		return FineVolumeSlideUp[TPeriod](cp)
	case 0xB: // Fine volume slide down
		return FineVolumeSlideDown[TPeriod](cp)
	case 0xC: // Note cut
		return NoteCut[TPeriod](cp)
	case 0xD: // Note delay
		return NoteDelay[TPeriod](cp)
	case 0xE: // Pattern delay
		return PatternDelay[TPeriod](cp)
	}
	return UnhandledCommand[TPeriod]{Command: ce, Info: cp}
}
