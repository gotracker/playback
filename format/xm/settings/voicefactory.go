package settings

import (
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVoice "github.com/gotracker/playback/format/xm/voice"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
)

type voiceFactory[TPeriod period.Period] struct{}

func (voiceFactory[TPeriod]) NewVoice() voice.RenderVoice[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning] {
	return xmVoice.New(GetMachineSettings[TPeriod]())
}

var (
	amigaVoiceFactory  voiceFactory[period.Amiga]
	linearVoiceFactory voiceFactory[period.Linear]
)