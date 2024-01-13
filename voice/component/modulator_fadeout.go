package component

import (
	"fmt"

	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
)

// FadeoutModulator is an amplitude (volume) modulator
type FadeoutModulator struct {
	settings FadeoutModulatorSettings
	active   bool
	vol      volume.Volume
}

type FadeoutModulatorSettings struct {
	GetEnabled func() bool
	Amount     volume.Volume
}

func (a *FadeoutModulator) Setup(settings FadeoutModulatorSettings) {
	a.settings = settings
	a.Reset()
}

func (a FadeoutModulator) Clone() FadeoutModulator {
	m := a
	return m
}

// Reset disables the fadeout and resets its volume
func (a *FadeoutModulator) Reset() {
	a.active = false
	a.vol = volume.Volume(1)
}

// Fadeout activates the fadeout
func (a *FadeoutModulator) Fadeout() {
	a.active = a.settings.Amount != 0
}

// SetActive sets the status of the fadeout active flag
func (a *FadeoutModulator) SetActive(active bool) {
	a.active = active
}

// IsEnabled returns the status of the fadeout enablement flag
func (a FadeoutModulator) IsActive() bool {
	if a.settings.GetEnabled == nil {
		return false
	}

	return a.settings.GetEnabled() && a.active
}

// GetVolume returns the value of the fadeout volume
func (a FadeoutModulator) GetVolume() volume.Volume {
	return a.vol
}

func (a FadeoutModulator) GetFinalVolume() volume.Volume {
	if !a.IsActive() {
		return volume.Volume(1)
	}
	return a.vol
}

// Advance advances the fadeout value by 1 tick
func (a *FadeoutModulator) Advance() {
	if a.IsActive() && a.vol > 0 {
		a.vol = min(max(a.vol-a.settings.Amount, 0), 1)
	}
}

func (a FadeoutModulator) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("active{%v} vol{%v}",
		a.active,
		a.vol,
	), comment)
}
