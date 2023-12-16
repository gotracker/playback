package effect

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
	"github.com/gotracker/playback/song"
)

type EffectS3M interface {
	playback.Effect
}

type ChannelCommand channel.DataEffect

// Factory produces an effect for the provided channel pattern data
func Factory(mem *channel.Memory, data song.ChannelData) EffectS3M {
	d, _ := data.(*channel.Data)
	if d == nil {
		return nil
	}

	if !d.What.HasCommand() {
		return nil
	}

	mem.LastNonZero(d.Info)
	switch d.Command + '@' {
	case '@': // unused
		return nil
	case 'A': // Set Speed
		return SetSpeed(d.Info)
	case 'B': // Pattern Jump
		return OrderJump(d.Info)
	case 'C': // Pattern Break
		return RowJump(d.Info)
	case 'D': // Volume Slide / Fine Volume Slide
		return volumeSlideFactory(mem, d.Command, d.Info)
	case 'E': // Porta Down/Fine Porta Down/Xtra Fine Porta
		xx := mem.LastNonZero(d.Info)
		x := xx >> 4
		if x == 0x0F {
			return FinePortaDown(xx)
		} else if x == 0x0E {
			return ExtraFinePortaDown(xx)
		}
		return PortaDown(d.Info)
	case 'F': // Porta Up/Fine Porta Up/Extra Fine Porta Down
		xx := mem.LastNonZero(d.Info)
		x := xx >> 4
		if x == 0x0F {
			return FinePortaUp(xx)
		} else if x == 0x0E {
			return ExtraFinePortaUp(xx)
		}
		return PortaUp(d.Info)
	case 'G': // Porta to note
		return PortaToNote(d.Info)
	case 'H': // Vibrato
		return Vibrato(d.Info)
	case 'I': // Tremor
		return Tremor(d.Info)
	case 'J': // Arpeggio
		return Arpeggio(d.Info)
	case 'K': // Vibrato+Volume Slide
		return NewVibratoVolumeSlide(mem, d.Command, d.Info)
	case 'L': // Porta+Volume Slide
		return NewPortaVolumeSlide(mem, d.Command, d.Info)
	case 'M': // unused
		return nil
	case 'N': // unused
		return nil
	case 'O': // Sample Offset
		return SampleOffset(d.Info)
	case 'P': // unused
		return nil
	case 'Q': // Retrig + Volume Slide
		return RetrigVolumeSlide(d.Info)
	case 'R': // Tremolo
		return Tremolo(d.Info)
	case 'S': // Special
		return specialEffect(mem, d)
	case 'T': // Set Tempo
		return SetTempo(d.Info)
	case 'U': // Fine Vibrato
		return FineVibrato(d.Info)
	case 'V': // Global Volume
		return SetGlobalVolume(d.Info)
	default:
	}
	return UnhandledCommand{Command: d.Command, Info: d.Info}
}

func specialEffect(mem *channel.Memory, data *channel.Data) EffectS3M {
	var cmd = mem.LastNonZero(data.Info)
	switch cmd >> 4 {
	case 0x0: // Set Filter on/off
		return EnableFilter(data.Info)
	//case 0x1: // Set Glissando on/off

	case 0x2: // Set FineTune
		return SetFinetune(data.Info)
	case 0x3: // Set Vibrato Waveform
		return SetVibratoWaveform(data.Info)
	case 0x4: // Set Tremolo Waveform
		return SetTremoloWaveform(data.Info)
	case 0x5: // unused
		return nil
	case 0x6: // Fine Pattern Delay
		return FinePatternDelay(data.Info)
	case 0x7: // unused
		return nil
	case 0x8: // Set Pan Position
		return SetPanPosition(data.Info)
	case 0xA: // Stereo Control
		return StereoControl(data.Info)
	case 0xB: // Pattern Loop
		return PatternLoop(data.Info)
	case 0xC: // Note Cut
		return NoteCut(data.Info)
	case 0xD: // Note Delay
		return NoteDelay(data.Info)
	case 0xE: // Pattern Delay
		return PatternDelay(data.Info)
	//case 0xF: // Funk Repeat (invert loop)
	default:
	}
	return UnhandledCommand{Command: data.Command, Info: data.Info}
}

func volumeSlideFactory(mem *channel.Memory, cd uint8, ce channel.DataEffect) EffectS3M {
	xy := mem.LastNonZero(ce)
	x := channel.DataEffect(xy >> 4)
	y := channel.DataEffect(xy & 0x0F)
	switch {
	case x == 0:
		return VolumeSlideDown(xy)
	case y == 0:
		return VolumeSlideUp(xy)
	case x == 0x0f:
		return FineVolumeSlideDown(xy)
	case y == 0x0f:
		return FineVolumeSlideUp(xy)
	}
	// There is a chance that a volume slide command is set with an invalid
	// value or is 00, in which case the memory might have the invalid value,
	// so we need to handle it by deferring to using a straight volume slide
	// down instead of panicking with an unhandled command, which mimics what
	// ScreamTracker 3.xx does.
	return VolumeSlideDown(xy)
	//return UnhandledCommand{Command: cd, Info: xy}
}
