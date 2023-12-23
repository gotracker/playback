package component

import (
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/envelope"
)

// VolumeEnvelope is an amplitude modulation envelope
type VolumeEnvelope struct {
	enabled   bool
	state     envelope.State[volume.Volume]
	vol       volume.Volume
	keyOn     bool
	prevKeyOn bool
}

func (e *VolumeEnvelope) Init(env *envelope.Envelope[volume.Volume]) {
	e.state.Init(env)
	e.Reset()
}

func (e VolumeEnvelope) Clone() VolumeEnvelope {
	return VolumeEnvelope{
		enabled:   e.enabled,
		state:     e.state.Clone(),
		vol:       e.vol,
		keyOn:     false,
		prevKeyOn: false,
	}
}

// Reset resets the state to defaults based on the envelope provided
func (e *VolumeEnvelope) Reset() {
	e.state.Reset()
	e.keyOn = false
	e.prevKeyOn = false
	e.update()
}

// SetEnabled sets the enabled flag for the envelope
func (e *VolumeEnvelope) SetEnabled(enabled bool) {
	e.enabled = enabled
}

// IsEnabled returns the enabled flag for the envelope
func (e *VolumeEnvelope) IsEnabled() bool {
	return e.enabled
}

// GetCurrentValue returns the current cached envelope value
func (e *VolumeEnvelope) GetCurrentValue() volume.Volume {
	return e.vol
}

// SetEnvelopePosition sets the current position in the envelope
func (e *VolumeEnvelope) SetEnvelopePosition(pos int) voice.Callback {
	keyOn := e.keyOn
	prevKeyOn := e.prevKeyOn
	e.state.Reset()
	// TODO: this is gross, but currently the most optimal way to find the correct position
	for i := 0; i < pos; i++ {
		if doneCB := e.Advance(keyOn, prevKeyOn); doneCB != nil {
			return doneCB
		}
	}
	return nil
}

// Advance advances the envelope state 1 tick and calculates the current envelope value
func (e *VolumeEnvelope) Advance(keyOn bool, prevKeyOn bool) voice.Callback {
	e.keyOn = keyOn
	e.prevKeyOn = prevKeyOn
	var doneCB voice.Callback
	if done := e.state.Advance(e.keyOn, e.prevKeyOn); done {
		doneCB = e.state.Envelope().OnFinished
	}
	e.update()
	return doneCB
}

func (e *VolumeEnvelope) update() {
	cur, next, t := e.state.GetCurrentValue(e.keyOn)

	var y0 volume.Volume
	if cur != nil {
		y0 = cur.Value()
	}

	var y1 volume.Volume
	if next != nil {
		y1 = next.Value()
	}

	e.vol = y0 + volume.Volume(t)*(y1-y0)
}
