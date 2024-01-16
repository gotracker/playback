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
	settings PitchPanModulatorSettings[TPanning]
	unkeyed  struct {
		enabled bool
		pitch   note.Semitone
	}
	keyed  struct{}
	panSep float32
}

type PitchPanModulatorSettings[TPanning types.Panning] struct {
	PitchPanEnable     bool
	PitchPanCenter     note.Semitone
	PitchPanSeparation float32
}

func (p *PitchPanModulator[TPanning]) Setup(settings PitchPanModulatorSettings[TPanning]) {
	p.settings = settings
	p.unkeyed.enabled = settings.PitchPanEnable
	p.unkeyed.pitch = settings.PitchPanCenter
	p.Reset()
}

func (p PitchPanModulator[TPanning]) Clone() PitchPanModulator[TPanning] {
	m := p
	return m
}

func (p *PitchPanModulator[TPanning]) Reset() error {
	return p.updatePitchPan()
}

// SetPitch updates the pan separation modulated by the provided pitch
func (p *PitchPanModulator[TPanning]) SetPitch(st note.Semitone) error {
	p.unkeyed.pitch = st
	return p.updatePitchPan()
}

// IsPitchPanEnabled returns the enablement of the pitch-pan separation function
func (p PitchPanModulator[TPanning]) IsPitchPanEnabled() bool {
	return p.unkeyed.enabled
}

// EnablePitchPan enables the pitch-pan separation function
func (p *PitchPanModulator[TPanning]) EnablePitchPan(enabled bool) error {
	p.unkeyed.enabled = enabled
	return p.updatePitchPan()
}

// SetPanSeparation gets the current pan separation
func (p PitchPanModulator[TPanning]) GetPanSeparation() float32 {
	return p.panSep
}

func (p PitchPanModulator[TPanning]) GetSeparatedPan(pan TPanning) TPanning {
	if !p.unkeyed.enabled || p.panSep == 0 {
		return pan
	}

	updatedPan := float32(pan) + p.panSep
	sepPan := TPanning(min(max(updatedPan, 0), float32(types.GetPanMax[TPanning]())))
	return sepPan
}

// Advance advances the fadeout value by 1 tick
func (p *PitchPanModulator[TPanning]) Advance() {
}

func (p *PitchPanModulator[TPanning]) updatePitchPan() error {
	if !p.unkeyed.enabled {
		return nil
	}

	p.panSep = (float32(p.unkeyed.pitch) - float32(p.settings.PitchPanCenter)) * p.settings.PitchPanSeparation
	return nil
}

func (p PitchPanModulator[TPanning]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("enabled{%v} pitch{%v} panSep{%v}",
		p.unkeyed.enabled,
		p.unkeyed.pitch,
		p.panSep,
	), comment)
}
