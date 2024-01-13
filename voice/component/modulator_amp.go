package component

import (
	"fmt"

	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/types"
)

// AmpModulator is an amplitude (volume) modulator
type AmpModulator[TMixingVolume, TVolume types.Volume] struct {
	settings AmpModulatorSettings[TMixingVolume, TVolume]
	active   bool
	vol      TVolume
	delta    types.VolumeDelta
	mixing   TMixingVolume
	final    volume.Volume // = active? * mixing * vol
}

type AmpModulatorSettings[TMixingVolume, TVolume types.Volume] struct {
	Active              bool
	DefaultMixingVolume TMixingVolume
	DefaultVolume       TVolume
}

func (a *AmpModulator[TMixingVolume, TVolume]) Setup(settings AmpModulatorSettings[TMixingVolume, TVolume]) {
	a.settings = settings
	a.active = settings.Active
	a.vol = settings.DefaultVolume
	a.delta = 0
	a.mixing = settings.DefaultMixingVolume
	a.updateFinal()
}

func (a AmpModulator[TMixingVolume, TVolume]) Clone() AmpModulator[TMixingVolume, TVolume] {
	m := a
	return m
}

func (a *AmpModulator[TMixingVolume, TVolume]) SetActive(active bool) {
	a.active = active
	a.updateFinal()
}

func (a AmpModulator[TMixingVolume, TVolume]) IsActive() bool {
	return a.active
}

// SetMixingVolume configures the mixing volume of the modulator
func (a *AmpModulator[TMixingVolume, TVolume]) SetMixingVolume(mixing TMixingVolume) {
	if !mixing.IsUseInstrumentVol() {
		a.mixing = mixing
		a.updateFinal()
	}
}

// GetMixingVolume returns the current mixing volume of the modulator
func (a AmpModulator[TMixingVolume, TVolume]) GetMixingVolume() TMixingVolume {
	return a.mixing
}

// SetVolume sets the current volume (before fadeout calculation)
func (a *AmpModulator[TMixingVolume, TVolume]) SetVolume(vol TVolume) {
	if vol.IsUseInstrumentVol() {
		vol = a.settings.DefaultVolume
	}
	a.vol = vol
	a.updateFinal()
}

// GetVolume returns the current volume (before fadeout calculation)
func (a AmpModulator[TMixingVolume, TVolume]) GetVolume() TVolume {
	return a.vol
}

func (a *AmpModulator[TMixingVolume, TVolume]) SetVolumeDelta(d types.VolumeDelta) {
	a.delta = d
	a.updateFinal()
}

func (a AmpModulator[TMixingVolume, TVolume]) GetVolumeDelta() types.VolumeDelta {
	return a.delta
}

// GetFinalVolume returns the current volume (after fadeout calculation)
func (a AmpModulator[TMixingVolume, TVolume]) GetFinalVolume() volume.Volume {
	return a.final
}

func (a *AmpModulator[TMixingVolume, TVolume]) updateFinal() {
	if !a.active {
		a.final = 0
		return
	}

	v := types.AddVolumeDelta(a.vol, a.delta)
	a.final = a.mixing.ToVolume() * v.ToVolume()
}

func (a AmpModulator[TMixingVolume, TVolume]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("active{%v} vol{%v} mixing{%v} final{%v}",
		a.active,
		a.vol,
		a.mixing,
		a.final,
	), comment)
}
