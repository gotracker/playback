package voice

import (
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"

	"github.com/heucuva/optional"
)

type envSettings struct {
	enabled optional.Value[bool]
	pos     optional.Value[int]
}

type playingMode uint8

const (
	playingModeAttack = playingMode(iota)
	playingModeRelease
)

type txn[TPeriod period.Period] struct {
	cancelled bool
	Voice     voice.Voice

	active      optional.Value[bool]
	playing     optional.Value[playingMode]
	fadeout     optional.Value[struct{}]
	period      optional.Value[TPeriod]
	periodDelta optional.Value[period.Delta]
	vol         optional.Value[volume.Volume]
	pos         optional.Value[sampling.Pos]
	pan         optional.Value[panning.Position]
	volEnv      envSettings
	pitchEnv    envSettings
	panEnv      envSettings
	filterEnv   envSettings
}

func (t *txn[TPeriod]) SetActive(active bool) {
	t.active.Set(active)
}

func (t *txn[TPeriod]) IsPendingActive() (bool, bool) {
	return t.active.Get()
}

func (t *txn[TPeriod]) IsCurrentlyActive() bool {
	return t.Voice.IsActive()
}

// Attack sets the playing mode to Attack
func (t *txn[TPeriod]) Attack() {
	t.playing.Set(playingModeAttack)
}

// Release sets the playing mode to Release
func (t *txn[TPeriod]) Release() {
	t.playing.Set(playingModeRelease)
}

// Fadeout activates the voice's fade-out function
func (t *txn[TPeriod]) Fadeout() {
	t.fadeout.Set(struct{}{})
}

// SetPeriod sets the period
func (t *txn[TPeriod]) SetPeriod(period TPeriod) {
	t.period.Set(period)
}

func (t *txn[TPeriod]) GetPendingPeriod() (TPeriod, bool) {
	return t.period.Get()
}

func (t *txn[TPeriod]) GetCurrentPeriod() TPeriod {
	return voice.GetPeriod[TPeriod](t.Voice)
}

// SetPeriodDelta sets the period delta
func (t *txn[TPeriod]) SetPeriodDelta(delta period.Delta) {
	t.periodDelta.Set(delta)
}

func (t *txn[TPeriod]) GetPendingPeriodDelta() (period.Delta, bool) {
	return t.periodDelta.Get()
}

func (t *txn[TPeriod]) GetCurrentPeriodDelta() period.Delta {
	return voice.GetPeriodDelta[TPeriod](t.Voice)
}

// SetVolume sets the volume
func (t *txn[TPeriod]) SetVolume(vol volume.Volume) {
	t.vol.Set(vol)
}

func (t *txn[TPeriod]) GetPendingVolume() (volume.Volume, bool) {
	return t.vol.Get()
}

func (t *txn[TPeriod]) GetCurrentVolume() volume.Volume {
	return voice.GetVolume(t.Voice)
}

// SetPos sets the position
func (t *txn[TPeriod]) SetPos(pos sampling.Pos) {
	t.pos.Set(pos)
}

func (t *txn[TPeriod]) GetPendingPos() (sampling.Pos, bool) {
	return t.pos.Get()
}

func (t *txn[TPeriod]) GetCurrentPos() sampling.Pos {
	return voice.GetPos(t.Voice)
}

// SetPan sets the panning position
func (t *txn[TPeriod]) SetPan(pan panning.Position) {
	t.pan.Set(pan)
}

func (t *txn[TPeriod]) GetPendingPan() (panning.Position, bool) {
	return t.pan.Get()
}

func (t *txn[TPeriod]) GetCurrentPan() panning.Position {
	return voice.GetPan(t.Voice)
}

// SetVolumeEnvelopePosition sets the volume envelope position
func (t *txn[TPeriod]) SetVolumeEnvelopePosition(pos int) {
	t.volEnv.pos.Set(pos)
}

// EnableVolumeEnvelope sets the volume envelope enable flag
func (t *txn[TPeriod]) EnableVolumeEnvelope(enabled bool) {
	t.volEnv.enabled.Set(enabled)
}

func (t *txn[TPeriod]) IsPendingVolumeEnvelopeEnabled() (bool, bool) {
	return t.volEnv.enabled.Get()
}

func (t *txn[TPeriod]) IsCurrentVolumeEnvelopeEnabled() bool {
	return voice.IsVolumeEnvelopeEnabled(t.Voice)
}

// SetPitchEnvelopePosition sets the pitch envelope position
func (t *txn[TPeriod]) SetPitchEnvelopePosition(pos int) {
	t.pitchEnv.pos.Set(pos)
}

// EnablePitchEnvelope sets the pitch envelope enable flag
func (t *txn[TPeriod]) EnablePitchEnvelope(enabled bool) {
	t.pitchEnv.enabled.Set(enabled)
}

// SetPanEnvelopePosition sets the panning envelope position
func (t *txn[TPeriod]) SetPanEnvelopePosition(pos int) {
	t.panEnv.pos.Set(pos)
}

// EnablePanEnvelope sets the pan envelope enable flag
func (t *txn[TPeriod]) EnablePanEnvelope(enabled bool) {
	t.panEnv.enabled.Set(enabled)
}

// SetFilterEnvelopePosition sets the pitch envelope position
func (t *txn[TPeriod]) SetFilterEnvelopePosition(pos int) {
	t.filterEnv.pos.Set(pos)
}

// EnableFilterEnvelope sets the filter envelope enable flag
func (t *txn[TPeriod]) EnableFilterEnvelope(enabled bool) {
	t.filterEnv.enabled.Set(enabled)
}

// SetAllEnvelopePositions sets all the envelope positions to the same value
func (t *txn[TPeriod]) SetAllEnvelopePositions(pos int) {
	t.volEnv.pos.Set(pos)
	t.pitchEnv.pos.Set(pos)
	t.panEnv.pos.Set(pos)
	t.filterEnv.pos.Set(pos)
}

// ======

// Cancel cancels a pending transaction
func (t *txn[TPeriod]) Cancel() {
	t.cancelled = true
}

// Commit commits the transaction by applying pending updates
func (t *txn[TPeriod]) Commit() {
	if t.cancelled {
		return
	}
	t.cancelled = true

	if t.Voice == nil {
		panic("voice not initialized")
	}

	if active, ok := t.active.Get(); ok {
		t.Voice.SetActive(active)
	}

	if p, ok := t.period.Get(); ok {
		voice.SetPeriod(t.Voice, p)
	}

	if delta, ok := t.periodDelta.Get(); ok {
		voice.SetPeriodDelta[TPeriod](t.Voice, delta)
	}

	if vol, ok := t.vol.Get(); ok {
		voice.SetVolume(t.Voice, vol)
	}

	if pos, ok := t.pos.Get(); ok {
		voice.SetPos(t.Voice, pos)
	}

	if pan, ok := t.pan.Get(); ok {
		voice.SetPan(t.Voice, pan)
	}

	if pos, ok := t.volEnv.pos.Get(); ok {
		voice.SetVolumeEnvelopePosition(t.Voice, pos)
	}

	if enabled, ok := t.volEnv.enabled.Get(); ok {
		voice.EnableVolumeEnvelope(t.Voice, enabled)
	}

	if pos, ok := t.pitchEnv.pos.Get(); ok {
		voice.SetPitchEnvelopePosition[TPeriod](t.Voice, pos)
	}

	if enabled, ok := t.pitchEnv.enabled.Get(); ok {
		voice.EnablePitchEnvelope[TPeriod](t.Voice, enabled)
	}

	if pos, ok := t.panEnv.pos.Get(); ok {
		voice.SetPanEnvelopePosition(t.Voice, pos)
	}

	if enabled, ok := t.panEnv.enabled.Get(); ok {
		voice.EnablePanEnvelope(t.Voice, enabled)
	}

	if pos, ok := t.filterEnv.pos.Get(); ok {
		voice.SetFilterEnvelopePosition(t.Voice, pos)
	}

	if enabled, ok := t.filterEnv.enabled.Get(); ok {
		voice.EnableFilterEnvelope(t.Voice, enabled)
	}

	if mode, ok := t.playing.Get(); ok {
		switch mode {
		case playingModeAttack:
			t.Voice.Attack()
		case playingModeRelease:
			t.Voice.Release()
		}
	}

	if _, ok := t.fadeout.Get(); ok {
		t.Voice.Fadeout()
	}
}

func (t *txn[TPeriod]) GetVoice() voice.Voice {
	return t.Voice
}

func (t *txn[TPeriod]) Clone() voice.Transaction[TPeriod] {
	c := *t
	return &c
}
