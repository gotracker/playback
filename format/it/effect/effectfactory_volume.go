package effect

import (
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

func volPanEffectFactory[TPeriod period.Period](mem *channel.Memory, v uint8) EffectIT {
	switch {
	case v <= 0x40: // volume set - handled elsewhere
		return nil
	case v >= 0x41 && v <= 0x4a: // fine volume slide up
		return VolChanFineVolumeSlideUp[TPeriod](v - 0x41)
	case v >= 0x4b && v <= 0x54: // fine volume slide down
		return VolChanFineVolumeSlideDown[TPeriod](v - 0x4b)
	case v >= 0x55 && v <= 0x5e: // volume slide up
		return VolChanVolumeSlideUp[TPeriod](v - 0x55)
	case v >= 0x5f && v <= 0x68: // volume slide down
		return VolChanVolumeSlideDown[TPeriod](v - 0x5f)
	case v >= 0x69 && v <= 0x72: // portamento down
		return volPortaDown[TPeriod](v - 0x69)
	case v >= 0x73 && v <= 0x7c: // portamento up
		return volPortaUp[TPeriod](v - 0x73)
	case v >= 0x80 && v <= 0xc0: // set panning
		return SetPanPosition[TPeriod](v - 0x80)
	case v >= 0xc1 && v <= 0xca: // portamento to note
		return volPortaToNote[TPeriod](v - 0xc1)
	case v >= 0xcb && v <= 0xd4: // vibrato
		return Vibrato[TPeriod](v - 0xcb)
	}
	return UnhandledVolCommand[TPeriod]{Vol: v}
}

func volPortaDown[TPeriod period.Period](v uint8) EffectIT {
	return PortaDown[TPeriod](v * 4)
}
func volPortaUp[TPeriod period.Period](v uint8) EffectIT {
	return PortaUp[TPeriod](v * 4)
}

func volPortaToNote[TPeriod period.Period](v uint8) EffectIT {
	switch v {
	case 0:
		return PortaToNote[TPeriod](0x00)
	case 1:
		return PortaToNote[TPeriod](0x01)
	case 2:
		return PortaToNote[TPeriod](0x04)
	case 3:
		return PortaToNote[TPeriod](0x08)
	case 4:
		return PortaToNote[TPeriod](0x10)
	case 5:
		return PortaToNote[TPeriod](0x20)
	case 6:
		return PortaToNote[TPeriod](0x40)
	case 7:
		return PortaToNote[TPeriod](0x60)
	case 8:
		return PortaToNote[TPeriod](0x80)
	case 9:
		return PortaToNote[TPeriod](0xFF)
	}
	// impossible, but hey...
	return UnhandledVolCommand[TPeriod]{Vol: v + 0xc1}
}
