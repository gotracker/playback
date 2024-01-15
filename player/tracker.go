package player

import (
	"errors"
	"time"

	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/sampler"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/tracing"
)

// Premixable is an interface to getting the premix data from the tracker
type Premixable interface {
	GetPremixData() (*output.PremixData, error)
}

// Tracker is an extensible music tracker
type Tracker struct {
	M machine.MachineTicker
	s *sampler.Sampler

	tracing.Tracing
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
	if t.M == nil {
		return nil, nil
	}

	defer t.OutputTraces()

	premix, err := t.M.Tick(t.s)
	if err != nil && !errors.Is(err, song.ErrStopSong) {
		return nil, err
	}

	return premix, err
}

// GetSampleRate returns the sample rate of the sampler
func (t *Tracker) GetSampleRate() frequency.Frequency {
	return frequency.Frequency(t.s.SampleRate)
}

// SetupSampler configures the internal sampler
func (t *Tracker) SetupSampler(samplesPerSecond int, channels int) error {
	t.s = sampler.NewSampler(samplesPerSecond, channels)
	if t.s == nil {
		return errors.New("NewSampler() returned nil")
	}

	return nil
}

// GetSampler returns the current sampler
func (t *Tracker) GetSampler() *sampler.Sampler {
	return t.s
}

func (t *Tracker) GetNumOrders() int {
	if t.M == nil {
		return 0
	}

	return t.M.GetNumOrders()
}

func (t *Tracker) CanOrderLoop() bool {
	if t.M == nil {
		return false
	}

	return t.M.CanOrderLoop()
}

func (t *Tracker) GetName() string {
	if t.M == nil {
		return ""
	}

	return t.M.GetName()
}

// Configure sets specified features
func (t *Tracker) Configure(features []feature.Feature) (settings.UserSettings, error) {
	us := settings.UserSettings{
		LongChannelOutput:    true,
		EnableNewNoteActions: true,
	}

	for _, feat := range features {
		switch f := feat.(type) {
		case feature.SongLoop:
			us.SongLoopCount = f.Count
		case feature.StartOrderAndRow:
			if o, set := f.Order.Get(); set {
				us.Start.Order.Set(index.Order(o))
			}
			if r, set := f.Row.Get(); set {
				us.Start.Row.Set(index.Row(r))
			}
		case feature.PlayUntilOrderAndRow:
			us.PlayUntil.Order.Set(index.Order(f.Order))
			us.PlayUntil.Row.Set(index.Row(f.Row))
		case feature.SetDefaultTempo:
			us.Start.Tempo = f.Tempo
		case feature.SetDefaultBPM:
			us.Start.BPM = f.BPM
		case feature.IgnoreUnknownEffect:
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
