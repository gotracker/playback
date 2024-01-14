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
	unkeyed  struct {
		enabled bool
	}
	keyed struct {
		vol volume.Volume
	}
}

type FadeoutModulatorSettings struct {
	Enabled   bool
	GetActive func() bool
	Amount    volume.Volume
}

func (a *FadeoutModulator) Setup(settings FadeoutModulatorSettings) {
	a.settings = settings
	a.unkeyed.enabled = settings.Enabled
	a.Reset()
}

func (a FadeoutModulator) Clone() FadeoutModulator {
	m := a
	return m
}

// Reset disables the fadeout and resets its volume
func (a *FadeoutModulator) Reset() {
	a.keyed.vol = volume.Volume(1)
}

// SetEnabled sets the status of the fadeout enable flag
func (a *FadeoutModulator) SetEnabled(enabled bool) {
	a.unkeyed.enabled = enabled
}

// IsEnabled returns the status of the fadeout enablement flag
func (a FadeoutModulator) IsActive() bool {
	if !a.unkeyed.enabled || a.settings.GetActive == nil {
		return false
	}

	return a.settings.GetActive()
}

// GetVolume returns the value of the fadeout volume
func (a FadeoutModulator) GetVolume() volume.Volume {
	return a.keyed.vol
}

func (a FadeoutModulator) GetFinalVolume() volume.Volume {
	if !a.IsActive() {
		return volume.Volume(1)
	}
	return a.keyed.vol
}

// Advance advances the fadeout value by 1 tick
func (a *FadeoutModulator) Advance() {
	if a.IsActive() && a.keyed.vol > 0 {
		a.keyed.vol = min(max(a.keyed.vol-a.settings.Amount, 0), 1)
	}
}

func (a FadeoutModulator) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("enabled{%v} vol{%v}",
		a.unkeyed.enabled,
		a.keyed.vol,
	), comment)
}
