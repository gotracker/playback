package voice

import (
	"github.com/gotracker/gomixing/panning"
	itPanning "github.com/gotracker/playback/format/it/panning"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/voice/types"
)

// == PanModulator ==

func (v *itVoice[TPeriod]) SetPan(pan itPanning.Panning) error {
	return v.pan.SetPan(pan)
}

func (v itVoice[TPeriod]) GetPan() itPanning.Panning {
	return v.pan.GetPan()
}

func (v *itVoice[TPeriod]) SetPanDelta(d types.PanDelta) error {
	return v.pan.SetPanDelta(d)
}

func (v itVoice[TPeriod]) GetPanDelta() types.PanDelta {
	return v.pan.GetPanDelta()
}

func (v itVoice[TPeriod]) GetPanSeparation() float32 {
	return v.pitchPan.GetPanSeparation()
}

func (v *itVoice[TPeriod]) SetPitchPanNote(st note.Semitone) error {
	return v.pitchPan.SetPitch(st)
}

func (v *itVoice[TPeriod]) EnablePitchPan(enabled bool) error {
	return v.pitchPan.EnablePitchPan(enabled)
}

func (v itVoice[TPeriod]) IsPitchPanEnabled() bool {
	return v.pitchPan.IsPitchPanEnabled()
}

func (v itVoice[TPeriod]) GetFinalPan() panning.Position {
	return v.finalPan
}
