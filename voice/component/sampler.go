package component

import (
	"fmt"

	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/loop"
	"github.com/gotracker/playback/voice/pcm"
	"github.com/gotracker/playback/voice/types"
)

// Sampler is a sampler component
type Sampler[TPeriod types.Period, TVolume types.Volume] struct {
	settings     SamplerSettings[TPeriod, TVolume]
	pos          sampling.Pos
	loopsEnabled bool

	slimKeyModulator
}

type SamplerSettings[TPeriod types.Period, TVolume types.Volume] struct {
	Sample        pcm.Sample
	DefaultVolume TVolume
	WholeLoop     loop.Loop
	SustainLoop   loop.Loop
}

func (s Sampler[TPeriod, TVolume]) Clone() Voicer[TPeriod, TVolume] {
	m := s
	return &m
}

// Setup sets up the sampler
func (s *Sampler[TPeriod, TVolume]) Setup(settings SamplerSettings[TPeriod, TVolume]) {
	s.settings = settings
}

// SetPos sets the current position of the sampler in the pcm data (and loops)
func (s *Sampler[TPeriod, TVolume]) SetPos(pos sampling.Pos) {
	s.pos = pos
}

// GetPos returns the current position of the sampler in the pcm data (and loops)
func (s Sampler[TPeriod, TVolume]) GetPos() sampling.Pos {
	return s.pos
}

// Attack sets the key-on value (for loop processing)
func (s *Sampler[TPeriod, TVolume]) Attack() {
	s.slimKeyModulator.Attack()
	s.loopsEnabled = true
}

// Release releases the key-on value (for loop processing)
func (s *Sampler[TPeriod, TVolume]) Release() {
	s.slimKeyModulator.Release()
}

// Fadeout disables the loops (for loop processing)
func (s *Sampler[TPeriod, TVolume]) Fadeout() {
	s.loopsEnabled = false
}

func (s *Sampler[TPeriod, TVolume]) DeferredAttack() {
	// does nothing
}

func (s *Sampler[TPeriod, TVolume]) DeferredRelease() {
	// does nothing
}

func (s Sampler[TPeriod, TVolume]) GetDefaultVolume() TVolume {
	return s.settings.DefaultVolume
}

func (s Sampler[TPeriod, TVolume]) GetNumChannels() int {
	if s.settings.Sample == nil {
		return 0
	}
	return s.settings.Sample.Channels()
}

// GetSample returns a multi-channel sample at the specified position
func (s *Sampler[TPeriod, TVolume]) GetSample(pos sampling.Pos) volume.Matrix {
	v0 := s.getConvertedSample(pos.Pos)
	if v0.Channels == 0 {
		if s.canLoop() {
			v01 := s.getConvertedSample(pos.Pos)
			panic(v01)
		}
		return v0
	}

	if pos.Frac == 0 {
		return v0
	}

	v1 := s.getConvertedSample(pos.Pos + 1)
	return v0.Lerp(v1, pos.Frac)
}

func (s Sampler[TPeriod, TVolume]) canLoop() bool {
	if s.loopsEnabled {
		return (s.keyOn && s.settings.SustainLoop.Enabled()) || s.settings.WholeLoop.Enabled()
	}
	return false
}

func (s *Sampler[TPeriod, TVolume]) getConvertedSample(pos int) volume.Matrix {
	if s.settings.Sample == nil {
		return volume.Matrix{}
	}
	sl := s.settings.Sample.Length()
	fadeout := false
	fadeoutLen := 0
	if pos >= sl {
		if s.canLoop() {
			pos, _ = loop.CalcLoopPos(s.settings.WholeLoop, s.settings.SustainLoop, pos, sl, s.keyOn)
		} else {
			fadeoutLen = pos - sl
			pos = sl - 1
			fadeout = true
		}
	}
	if pos < 0 || pos >= sl {
		return volume.Matrix{}
	}
	s.settings.Sample.Seek(pos)
	data, err := s.settings.Sample.Read()
	if err != nil {
		return volume.Matrix{}
	}

	if !fadeout {
		return data
	}
	if fadeoutLen >= 32 {
		return data.Apply(0)
	}

	atten := volume.Volume(1) / volume.Volume(int(1<<fadeoutLen))
	return data.Apply(atten)
}

func (s Sampler[TPeriod, TVolume]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("pos{%v} loopsEnabled{%v}",
		s.pos,
		s.loopsEnabled,
	), comment)
}
