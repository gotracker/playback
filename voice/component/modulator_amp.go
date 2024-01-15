package component

import (
	"fmt"

	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/types"
	"github.com/heucuva/optional"
)

// AmpModulator is an amplitude (volume) modulator
type AmpModulator[TMixingVolume, TVolume types.Volume] struct {
	settings AmpModulatorSettings[TMixingVolume, TVolume]

	unkeyed struct {
		active bool
		vol    TVolume
		mixing TMixingVolume
	}
	keyed struct {
		delta          types.VolumeDelta
		mixingOverride optional.Value[TMixingVolume]
	}
	final volume.Volume // = active? * mixing * vol
}

type AmpModulatorSettings[TMixingVolume, TVolume types.Volume] struct {
	Active              bool
	DefaultMixingVolume TMixingVolume
	DefaultVolume       TVolume
}

func (a *AmpModulator[TMixingVolume, TVolume]) Setup(settings AmpModulatorSettings[TMixingVolume, TVolume]) {
	a.settings = settings
	a.unkeyed.active = settings.Active
	a.unkeyed.vol = settings.DefaultVolume
	a.unkeyed.mixing = settings.DefaultMixingVolume
	a.Reset()
}

func (a AmpModulator[TMixingVolume, TVolume]) Clone() AmpModulator[TMixingVolume, TVolume] {
	m := a
	return m
}

func (a *AmpModulator[TMixingVolume, TVolume]) Reset() {
	a.keyed.delta = 0
	a.keyed.mixingOverride.Reset()
	a.updateFinal()
}

func (a *AmpModulator[TMixingVolume, TVolume]) SetActive(active bool) {
	a.unkeyed.active = active
	a.updateFinal()
}

func (a AmpModulator[TMixingVolume, TVolume]) IsActive() bool {
	return a.unkeyed.active
}

// SetMixingVolume configures the mixing volume of the modulator
func (a *AmpModulator[TMixingVolume, TVolume]) SetMixingVolume(mixing TMixingVolume) {
	if !mixing.IsUseInstrumentVol() {
		a.unkeyed.mixing = mixing
		a.updateFinal()
	}
}

// GetMixingVolume returns the current mixing volume of the modulator
func (a AmpModulator[TMixingVolume, TVolume]) GetMixingVolume() TMixingVolume {
	return a.unkeyed.mixing
}

func (a *AmpModulator[TMixingVolume, TVolume]) SetMixingVolumeOverride(mvo optional.Value[TMixingVolume]) {
	a.keyed.mixingOverride = mvo
}

func (a AmpModulator[TMixingVolume, TVolume]) GetMixingVolumeOverride() optional.Value[TMixingVolume] {
	return a.keyed.mixingOverride
}

// SetVolume sets the current volume (before fadeout calculation)
func (a *AmpModulator[TMixingVolume, TVolume]) SetVolume(vol TVolume) {
	if vol.IsUseInstrumentVol() {
		vol = a.settings.DefaultVolume
	}
	a.unkeyed.vol = vol
	a.updateFinal()
}

// GetVolume returns the current volume (before fadeout calculation)
func (a AmpModulator[TMixingVolume, TVolume]) GetVolume() TVolume {
	return a.unkeyed.vol
}

func (a *AmpModulator[TMixingVolume, TVolume]) SetVolumeDelta(d types.VolumeDelta) {
	a.keyed.delta = d
	a.updateFinal()
}

func (a AmpModulator[TMixingVolume, TVolume]) GetVolumeDelta() types.VolumeDelta {
	return a.keyed.delta
}

// GetFinalVolume returns the current volume (after fadeout calculation)
func (a AmpModulator[TMixingVolume, TVolume]) GetFinalVolume() volume.Volume {
	return a.final
}

func (a *AmpModulator[TMixingVolume, TVolume]) updateFinal() {
	if !a.unkeyed.active {
		a.final = 0
		return
	}

	v := types.AddVolumeDelta(a.unkeyed.vol, a.keyed.delta)

	mv := a.unkeyed.mixing
	if mvo, set := a.keyed.mixingOverride.Get(); set {
		mv = mvo
	}

	a.final = mv.ToVolume() * v.ToVolume()
}

func (a AmpModulator[TMixingVolume, TVolume]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("active{%v} vol{%v} mixing{%v} mixingOverride{%v} delta{%v} final{%v}",
		a.unkeyed.active,
		a.unkeyed.vol,
		a.unkeyed.mixing,
		a.keyed.mixingOverride,
		a.keyed.delta,
		a.final,
	), comment)
}
