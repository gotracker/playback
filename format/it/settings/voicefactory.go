package settings

import (
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVoice "github.com/gotracker/playback/format/it/voice"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
)

type voiceFactory[TPeriod period.Period] struct{}

func (voiceFactory[TPeriod]) NewVoice() voice.RenderVoice[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning] {
	return itVoice.New(GetMachineSettings[TPeriod]())
}

var (
	amigaVoiceFactory  voiceFactory[period.Amiga]
	linearVoiceFactory voiceFactory[period.Linear]
)
