package channel

import (
	"github.com/gotracker/playback"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/song"
)

type ChannelCommand DataEffect

// Factory produces an effect for the provided channel pattern data
func EffectFactory(mem *Memory, data song.ChannelData[s3mVolume.Volume]) playback.Effect {
	if data == nil {
		return nil
	}

	d, _ := data.(Data)
	if !d.What.HasCommand() {
		return nil
	}

	// Store the last non-zero value
	_ = mem.LastNonZero(d.Info)

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
		xx := mem.Porta(d.Info)
		x := xx >> 4
		if x == 0x0F {
			return FinePortaDown(xx)
		} else if x == 0x0E {
			return ExtraFinePortaDown(xx)
		}
		return PortaDown(d.Info)
	case 'F': // Porta Up/Fine Porta Up/Extra Fine Porta Down
		xx := mem.Porta(d.Info)
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

func specialEffect(mem *Memory, data Data) playback.Effect {
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
	case 0x9: // Sound Control
		return soundControlEffect(data)
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

func volumeSlideFactory(mem *Memory, cd uint8, ce DataEffect) playback.Effect {
	xy := mem.LastNonZero(ce)
	x := DataEffect(xy >> 4)
	y := DataEffect(xy & 0x0F)
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

func soundControlEffect(data Data) playback.Effect {
	switch data.Info & 0xF {
	case 0x0: // Surround Off
	case 0x1: // Surround On
		// only S91 is supported directly by S3M
		return SurroundOn(data.Info)
	case 0x8: // Reverb Off
	case 0x9: // Reverb On
	case 0xA: // Center Surround
	case 0xB: // Quad Surround
	case 0xC: // Global Filters
	case 0xD: // Local Filters
	case 0xE: // Play Forward
	case 0xF: // Play Backward
	}
	return UnhandledCommand{Command: data.Command, Info: data.Info}
}
