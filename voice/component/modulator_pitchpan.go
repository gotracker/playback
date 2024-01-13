package component

import (
	"fmt"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/types"
)

// PitchPanModulator is an pan (spatial) modulator
type PitchPanModulator[TPanning types.Panning] struct {
	settings        PitchPanModulatorSettings[TPanning]
	pitchPanEnabled bool
	pitch           note.Semitone
	panSep          float32
}

type PitchPanModulatorSettings[TPanning types.Panning] struct {
	PitchPanEnable     bool
	PitchPanCenter     note.Semitone
	PitchPanSeparation float32
}

func (p *PitchPanModulator[TPanning]) Setup(settings PitchPanModulatorSettings[TPanning]) {
	p.settings = settings
	p.pitchPanEnabled = settings.PitchPanEnable
	p.pitch = settings.PitchPanCenter
	p.Reset()
}

func (p PitchPanModulator[TPanning]) Clone() PitchPanModulator[TPanning] {
	m := p
	return m
}

func (p *PitchPanModulator[TPanning]) Reset() {
	p.updatePitchPan()
}

// SetPitch updates the pan separation modulated by the provided pitch
func (p *PitchPanModulator[TPanning]) SetPitch(st note.Semitone) {
	p.pitch = st
	p.updatePitchPan()
}

// IsPitchPanEnabled returns the enablement of the pitch-pan separation function
func (p PitchPanModulator[TPanning]) IsPitchPanEnabled() bool {
	return p.pitchPanEnabled
}

// EnablePitchPan enables the pitch-pan separation function
func (p *PitchPanModulator[TPanning]) EnablePitchPan(enabled bool) {
	p.pitchPanEnabled = enabled
	p.updatePitchPan()
}

// SetPanSeparation gets the current pan separation
func (p PitchPanModulator[TPanning]) GetPanSeparation() float32 {
	return p.panSep
}

func (p PitchPanModulator[TPanning]) GetSeparatedPan(pan TPanning) TPanning {
	if !p.pitchPanEnabled || p.panSep == 0 {
		return pan
	}

	updatedPan := float32(pan) + p.panSep
	sepPan := TPanning(min(max(updatedPan, 0), float32(types.GetPanMax[TPanning]())))
	return sepPan
}

// Advance advances the fadeout value by 1 tick
func (p *PitchPanModulator[TPanning]) Advance() {
}

func (p *PitchPanModulator[TPanning]) updatePitchPan() {
	if !p.pitchPanEnabled {
		return
	}

	p.panSep = (float32(p.pitch) - float32(p.settings.PitchPanCenter)) * p.settings.PitchPanSeparation
}

func (p PitchPanModulator[TPanning]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("pitchPanEnabled{%v} pitch{%v} panSep{%v}",
		p.pitchPanEnabled,
		p.pitch,
		p.panSep,
	), comment)
}
