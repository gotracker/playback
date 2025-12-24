package instrument

import (
	"github.com/heucuva/optional"

	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/voice/envelope"
	"github.com/gotracker/playback/voice/fadeout"
	"github.com/gotracker/playback/voice/loop"
	"github.com/gotracker/playback/voice/pcm"
	"github.com/gotracker/playback/voice/pitchpan"
	"github.com/gotracker/playback/voice/types"
)

// PCM is a PCM-data instrument
type PCM[TMixingVolume, TVolume types.Volume, TPanning types.Panning] struct {
	Sample               pcm.Sample
	Loop                 loop.Loop
	SustainLoop          loop.Loop
	Panning              optional.Value[TPanning]
	MixingVolume         optional.Value[TMixingVolume]
	FadeOut              fadeout.Settings
	PitchPan             pitchpan.PitchPan
	VolEnv               envelope.Envelope[TVolume]
	VolEnvFinishFadesOut bool
	PanEnv               envelope.Envelope[TPanning]
	PitchFiltMode        bool                                    // true = filter, false = pitch
	PitchFiltEnv         envelope.Envelope[types.PitchFiltValue] // this is either pitch or filter
}

func (p PCM[TMixingVolume, TVolume, TPanning]) GetLength() sampling.Pos {
	return sampling.Pos{Pos: p.Sample.Length()}
}
