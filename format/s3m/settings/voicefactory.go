package settings

import (
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVoice "github.com/gotracker/playback/format/s3m/voice"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
)

type voiceFactory struct{}

func (voiceFactory) NewVoice() voice.RenderVoice[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning] {
	return s3mVoice.New(GetMachineSettings())
}

var (
	amigaVoiceFactory voiceFactory
)
