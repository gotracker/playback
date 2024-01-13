package player

import (
	"errors"
	"time"

	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/player/sampler"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/tracing"
	voiceRender "github.com/gotracker/playback/voice/render"
)

// Premixable is an interface to getting the premix data from the tracker
type Premixable interface {
	GetPremixData() (*output.PremixData, error)
}

// Tracker is an extensible music tracker
type Tracker struct {
	BaseClockRate period.Frequency
	Premixable    Premixable

	s    *sampler.Sampler
	opl2 voiceRender.OPL2Chip

	globalVolume volume.Volume
	mixerVolume  volume.Volume

	ignoreUnknownEffect feature.IgnoreUnknownEffect
	outputChannels      map[int]render.ChannelIntf

	tracing.Tracing

	M                 machine.MachineTicker
	LongChannelOutput bool
}

func (t *Tracker) SetupMachine(s song.Data, us settings.UserSettings) error {
	var err error
	t.M, err = machine.NewMachine(s, us)
	return err
}

func (t *Tracker) Close() {
	t.Trace("Close")
	t.Tracing.Close()
}

// Update runs processing on the tracker, producing premixed sound data
func (t *Tracker) Update(deltaTime time.Duration, out chan<- *output.PremixData) error {
	premix, err := t.Generate(deltaTime)
	if premix != nil && len(premix.Data) > 0 {
		out <- premix
	}

	return err
}

func (t *Tracker) Generate(deltaTime time.Duration) (*output.PremixData, error) {
	defer t.OutputTraces()

	var (
		premix *output.PremixData
		err    error
	)
	if t.M != nil {
		premix, err = t.M.Tick(t.s)
		if err != nil && !errors.Is(err, song.ErrStopSong) {
			return nil, err
		}
	}

	if t.opl2 != nil {
		rr := [1]mixing.Data{}
		t.renderOPL2Tick(&rr[0],
			t.s.Mixer(),
			premix.SamplesLen)
		premix.Data = append(premix.Data, rr[:])

		// make room in the mixer for the OPL2 data
		// effectively, we can do this by calculating the new number (+1) of channels from the mixer volume (channels = reciprocal of mixer volume):
		//   numChannels = (1/mv) + 1
		// then by taking the reciprocal of it:
		//   1 / numChannels
		// but that ends up being simplified to:
		//   mv / (mv + 1)
		// and we get protection from div/0 in the process - provided, of course, that the mixerVolume is not exactly -1...
		mv := premix.MixerVolume
		premix.MixerVolume /= (mv + 1)
	}

	return premix, err
}

// GetRenderChannel returns the output channel for the provided index `ch`
func (t *Tracker) GetRenderChannel(ch int, init func(ch int) render.ChannelIntf) render.ChannelIntf {
	if t.outputChannels == nil {
		t.outputChannels = make(map[int]render.ChannelIntf)
	}

	if oc, ok := t.outputChannels[ch]; ok {
		return oc
	}
	oc := init(ch)
	t.outputChannels[ch] = oc
	return oc
}

// GetSampleRate returns the sample rate of the sampler
func (t *Tracker) GetSampleRate() period.Frequency {
	return period.Frequency(t.GetSampler().SampleRate)
}

func (t *Tracker) renderOPL2Tick(mixerData *mixing.Data, mix *mixing.Mixer, tickSamples int) {
	// make a stand-alone data buffer for this channel for this tick
	data := mix.NewMixBuffer(tickSamples)

	opl2data := make([]int32, tickSamples)

	if opl2 := t.opl2; opl2 != nil {
		opl2.GenerateBlock2(uint(tickSamples), opl2data)
	}

	for i, s := range opl2data {
		sv := volume.Volume(s) / 32768.0
		data[i].Assign(1, []volume.Volume{sv})
	}
	*mixerData = mixing.Data{
		Data:       data,
		Pan:        panning.CenterAhead,
		Volume:     t.globalVolume,
		SamplesLen: tickSamples,
	}
}

// GetOPL2Chip returns the current song's OPL2 chip, if it's needed
func (t *Tracker) GetOPL2Chip() voiceRender.OPL2Chip {
	return t.opl2
}

// SetOPL2Chip sets the current song's OPL2 chip
func (t *Tracker) SetOPL2Chip(opl2 voiceRender.OPL2Chip) {
	t.opl2 = opl2
}

// SetupSampler configures the internal sampler
func (t *Tracker) SetupSampler(samplesPerSecond int, channels int) error {
	t.s = sampler.NewSampler(samplesPerSecond, channels, t.BaseClockRate)
	if t.s == nil {
		return errors.New("NewSampler() returned nil")
	}

	return nil
}

// GetSampler returns the current sampler
func (t *Tracker) GetSampler() *sampler.Sampler {
	return t.s
}

// GetGlobalVolume returns the global volume value
func (t *Tracker) GetGlobalVolume() volume.Volume {
	return t.globalVolume
}

// SetGlobalVolume sets the global volume to the specified `vol` value
func (t *Tracker) SetGlobalVolume(vol volume.Volume) {
	t.TraceValueChange("SetGlobalVolume", t.globalVolume, vol)
	t.globalVolume = vol
}

// GetMixerVolume returns the mixer volume value
func (t *Tracker) GetMixerVolume() volume.Volume {
	return t.mixerVolume
}

// SetMixerVolume sets the mixer volume to the specified `vol` value
func (t *Tracker) SetMixerVolume(vol volume.Volume) {
	t.TraceValueChange("SetMixerVolume", t.mixerVolume, vol)
	t.mixerVolume = vol
}

// IgnoreUnknownEffect returns true if the tracker wants unknown effects to be ignored
func (t *Tracker) IgnoreUnknownEffect() bool {
	return t.ignoreUnknownEffect.Enabled
}

// Configure sets specified features
func (t *Tracker) Configure(features []feature.Feature) (settings.UserSettings, error) {
	us := settings.UserSettings{
		LongChannelOutput:    true,
		EnableNewNoteActions: true,
	}

	for _, feat := range features {
		switch f := feat.(type) {
		case feature.IgnoreUnknownEffect:
			t.ignoreUnknownEffect = f
			us.IgnoreUnknownEffect = f.Enabled
		case feature.EnableTracing:
			if err := t.EnableTracing(f.Filename); err != nil {
				return us, err
			}
			us.Tracer = &t.Tracing
		}
	}
	return us, nil
}
