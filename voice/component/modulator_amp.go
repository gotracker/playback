package component

import (
	"fmt"

	"github.com/heucuva/optional"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/types"
)

// AmpModulator is an amplitude (volume) modulator
type AmpModulator[TMixingVolume, TVolume types.Volume] struct {
	settings AmpModulatorSettings[TMixingVolume, TVolume]

	unkeyed struct {
		active bool
		muted  bool
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
	Muted               bool
	DefaultMixingVolume TMixingVolume
	DefaultVolume       TVolume
}

func (a *AmpModulator[TMixingVolume, TVolume]) Setup(settings AmpModulatorSettings[TMixingVolume, TVolume]) {
	a.settings = settings
	a.unkeyed.active = settings.Active
	a.unkeyed.muted = settings.Muted
	a.unkeyed.vol = settings.DefaultVolume
	a.unkeyed.mixing = settings.DefaultMixingVolume
	a.Reset()
}

func (a AmpModulator[TMixingVolume, TVolume]) Clone() AmpModulator[TMixingVolume, TVolume] {
	m := a
	return m
}

func (a *AmpModulator[TMixingVolume, TVolume]) Reset() error {
	a.keyed.delta = 0
	a.keyed.mixingOverride.Reset()
	return a.updateFinal()
}

func (a *AmpModulator[TMixingVolume, TVolume]) SetActive(active bool) error {
	a.unkeyed.active = active
	return a.updateFinal()
}

func (a AmpModulator[TMixingVolume, TVolume]) IsActive() bool {
	return a.unkeyed.active
}

func (a *AmpModulator[TMixingVolume, TVolume]) SetMuted(muted bool) error {
	a.unkeyed.muted = muted
	return nil
}

func (a AmpModulator[TMixingVolume, TVolume]) IsMuted() bool {
	return a.unkeyed.muted
}

// SetMixingVolume configures the mixing volume of the modulator
func (a *AmpModulator[TMixingVolume, TVolume]) SetMixingVolume(mixing TMixingVolume) error {
	if mixing.IsUseInstrumentVol() {
		return nil
	}

	a.unkeyed.mixing = mixing
	return a.updateFinal()
}

// GetMixingVolume returns the current mixing volume of the modulator
func (a AmpModulator[TMixingVolume, TVolume]) GetMixingVolume() TMixingVolume {
	return a.unkeyed.mixing
}

func (a *AmpModulator[TMixingVolume, TVolume]) SetMixingVolumeOverride(mvo optional.Value[TMixingVolume]) error {
	a.keyed.mixingOverride = mvo
	return nil
}

func (a AmpModulator[TMixingVolume, TVolume]) GetMixingVolumeOverride() optional.Value[TMixingVolume] {
	return a.keyed.mixingOverride
}

// SetVolume sets the current volume (before fadeout calculation)
func (a *AmpModulator[TMixingVolume, TVolume]) SetVolume(vol TVolume) error {
	if vol.IsUseInstrumentVol() {
		vol = a.settings.DefaultVolume
	}
	a.unkeyed.vol = vol
	return a.updateFinal()
}

// GetVolume returns the current volume (before fadeout calculation)
func (a AmpModulator[TMixingVolume, TVolume]) GetVolume() TVolume {
	return a.unkeyed.vol
}

func (a *AmpModulator[TMixingVolume, TVolume]) SetVolumeDelta(d types.VolumeDelta) error {
	a.keyed.delta = d
	return a.updateFinal()
}

func (a AmpModulator[TMixingVolume, TVolume]) GetVolumeDelta() types.VolumeDelta {
	return a.keyed.delta
}

// GetFinalVolume returns the current volume (after fadeout calculation)
func (a AmpModulator[TMixingVolume, TVolume]) GetFinalVolume() volume.Volume {
	return a.final
}

func (a *AmpModulator[TMixingVolume, TVolume]) updateFinal() error {
	if !a.unkeyed.active {
		a.final = 0
		return nil
	}

	v := types.AddVolumeDelta(a.unkeyed.vol, a.keyed.delta)

	mv := a.unkeyed.mixing
	if mvo, set := a.keyed.mixingOverride.Get(); set {
		mv = mvo
	}

	a.final = mv.ToVolume() * v.ToVolume()
	return nil
}

func (a AmpModulator[TMixingVolume, TVolume]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("active{%v} muted{%v} vol{%v} mixing{%v} mixingOverride{%v} delta{%v} final{%v}",
		a.unkeyed.active,
		a.unkeyed.muted,
		a.unkeyed.vol,
		a.unkeyed.mixing,
		a.keyed.mixingOverride,
		a.keyed.delta,
		a.final,
	), comment)
}
