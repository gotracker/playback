package state

import (
	"time"

	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
)

// Active is the active state of a channel
type Active[TPeriod period.Period] struct {
	playback.ChannelState[TPeriod]
	Voice       voice.Voice
	PeriodDelta period.Delta
}

// Reset sets the active state to defaults
func (a *Active[TPeriod]) Reset() {
	if v := a.Voice; v != nil {
		v.Release()
		a.Voice = nil
	}
	a.PeriodDelta = 0
	a.ChannelState.Reset()
}

// Clone clones the active state so that various interfaces do not collide
func (a *Active[TPeriod]) Clone() *Active[TPeriod] {
	var c Active[TPeriod] = *a
	if a.Voice != nil {
		c.Voice = a.Voice.Clone()
	}

	return &c
}

type RenderDetails struct {
	Mix          *mixing.Mixer
	Panmixer     mixing.PanMixer
	SamplerSpeed float32
	Samples      int
	Duration     time.Duration
}

// RenderStatesTogether renders a channel's series of sample data for a the provided number of samples
func RenderStatesTogether[TPeriod period.Period](periodConverter period.PeriodConverter[TPeriod], activeState *Active[TPeriod], pastNotes []*Active[TPeriod], details RenderDetails) []mixing.Data {
	var mixData []mixing.Data

	centerAheadPan := details.Panmixer.GetMixingMatrix(panning.CenterAhead)

	if activeState != nil {
		if data := activeState.renderState(periodConverter, centerAheadPan, details); data != nil {
			mixData = append(mixData, *data)
		}
	}

	for _, pn := range pastNotes {
		if pn != nil {
			if data := pn.renderState(periodConverter, centerAheadPan, details); data != nil {
				mixData = append(mixData, *data)
			}
		}
	}

	return mixData
}

func (a *Active[TPeriod]) renderState(periodConverter period.PeriodConverter[TPeriod], centerAheadPan volume.Matrix, details RenderDetails) *mixing.Data {
	if a.Period.IsInvalid() {
		return nil
	}

	vol := a.GetVolume()
	if vol == 0 {
		return nil
	}

	ncv := a.Voice
	if ncv == nil || ncv.IsDone() {
		return nil
	}

	// Commit the playback settings to the note-control
	voice.SetPeriod(ncv, a.Period)
	voice.SetVolume(ncv, vol)
	voice.SetPos(ncv, a.Pos)
	voice.SetPan(ncv, a.Pan)

	voice.SetPeriodDelta[TPeriod](ncv, a.PeriodDelta)

	// the period might be updated by the auto-vibrato system, here
	ncv.Advance(details.Duration)

	if !ncv.IsActive() {
		return nil
	}

	sampler := ncv.GetSampler(details.SamplerSpeed)

	if sampler == nil {
		return nil
	}

	// ... so grab the new value now.
	p := voice.GetFinalPeriod[TPeriod](ncv)
	pan := voice.GetFinalPan(ncv)

	// make a stand-alone data buffer for this channel for this tick
	sampleData := mixing.SampleMixIn{
		Sample:    sampler,
		StaticVol: volume.Volume(1.0),
		VolMatrix: centerAheadPan,
		MixPos:    0,
		MixLen:    details.Samples,
	}

	mixBuffer := details.Mix.NewMixBuffer(details.Samples)
	mixBuffer.MixInSample(sampleData)
	data := &mixing.Data{
		Data:       mixBuffer,
		Pan:        pan,
		Volume:     volume.Volume(1.0),
		Pos:        0,
		SamplesLen: details.Samples,
	}

	a.Pos = voice.GetPos(ncv)
	samplerAdd := float32(periodConverter.GetSamplerAdd(p, float64(details.SamplerSpeed)))
	a.Pos.Add(samplerAdd * float32(details.Samples))

	return data
}
