package voice

import (
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/voice/types"
)

// == PanModulator ==

func (v *xmVoice[TPeriod]) SetPan(pan xmPanning.Panning) error {
	return v.pan.SetPan(pan)
}

func (v xmVoice[TPeriod]) GetPan() xmPanning.Panning {
	return v.pan.GetPan()
}

func (v *xmVoice[TPeriod]) SetPanDelta(d types.PanDelta) error {
	return v.pan.SetPanDelta(d)
}

func (v xmVoice[TPeriod]) GetPanDelta() types.PanDelta {
	return v.pan.GetPanDelta()
}

func (v xmVoice[TPeriod]) GetFinalPan() panning.Position {
	if !v.IsPanEnvelopeEnabled() {
		return v.pan.GetFinalPan()
	}

	envPan := v.panEnv.GetCurrentValue()
	return envPan.ToPosition()
}
