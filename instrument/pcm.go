package instrument

import (
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/voice/envelope"
	"github.com/gotracker/playback/voice/fadeout"
	"github.com/gotracker/playback/voice/loop"
	"github.com/gotracker/playback/voice/pcm"
)

// PCM is a PCM-data instrument
type PCM struct {
	Sample        pcm.Sample
	Loop          loop.Loop
	SustainLoop   loop.Loop
	Panning       panning.Position
	MixingVolume  volume.Volume
	FadeOut       fadeout.Settings
	VolEnv        envelope.Envelope[volume.Volume]
	PanEnv        envelope.Envelope[panning.Position]
	PitchFiltMode bool                    // true = filter, false = pitch
	PitchFiltEnv  envelope.Envelope[int8] // this is either pitch or filter
}

func (PCM) GetKind() Kind {
	return KindPCM
}

func (p PCM) GetLength() sampling.Pos {
	return sampling.Pos{Pos: p.Sample.Length()}
}
