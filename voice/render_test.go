package voice

import (
	"testing"

	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/mixing"
	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/system"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/mixer"
)

type stubPeriod float32

func (stubPeriod) IsInvalid() bool { return false }

type stubSystem struct{}

func (stubSystem) GetMaxPastNotesPerChannel() int     { return 0 }
func (stubSystem) GetCommonRate() frequency.Frequency { return 0 }

type stubPeriodConverter struct {
	samplerAdd float64
}

func (s stubPeriodConverter) GetSystem() system.System       { return stubSystem{} }
func (s stubPeriodConverter) GetPeriod(note.Note) stubPeriod { return 0 }
func (s stubPeriodConverter) PortaToNote(p stubPeriod, d period.Delta, target stubPeriod) (stubPeriod, error) {
	return p, nil
}
func (s stubPeriodConverter) PortaDown(p stubPeriod, d period.Delta) (stubPeriod, error) {
	return p, nil
}
func (s stubPeriodConverter) PortaUp(p stubPeriod, d period.Delta) (stubPeriod, error) { return p, nil }
func (s stubPeriodConverter) AddDelta(p stubPeriod, d period.Delta) (stubPeriod, error) {
	return p, nil
}
func (s stubPeriodConverter) GetSamplerAdd(p stubPeriod, _ frequency.Frequency, _ frequency.Frequency) float64 {
	return s.samplerAdd
}
func (s stubPeriodConverter) GetFrequency(p stubPeriod) frequency.Frequency {
	return frequency.Frequency(p)
}

type stubPanMatrix struct {
	factorL float32
	factorR float32
	calls   int
}

func (p *stubPanMatrix) ApplyToMatrix(mtx volume.Matrix) volume.Matrix {
	p.calls++
	out := mtx
	if out.Channels >= 1 {
		out.StaticMatrix[0] *= volume.Volume(p.factorL)
	}
	if out.Channels >= 2 {
		out.StaticMatrix[1] *= volume.Volume(p.factorR)
	}
	return out
}

func (p *stubPanMatrix) Apply(v volume.Volume) volume.Matrix {
	p.calls++
	return volume.Matrix{StaticMatrix: volume.StaticMatrix{v * volume.Volume(p.factorL), v * volume.Volume(p.factorR)}, Channels: 2}
}

type stubDetailPanMixer struct {
	matrix panning.PanMixer
	pan    panning.Position
	sep    float32
	calls  int
}

func (p *stubDetailPanMixer) GetMixingMatrix(pan panning.Position, sep float32) panning.PanMixer {
	p.calls++
	p.pan = pan
	p.sep = sep
	return p.matrix
}

func (p *stubDetailPanMixer) NumChannels() int { return 2 }

type stubFilter struct {
	calls int
}

func (f *stubFilter) ApplyFilter(dry volume.Matrix) volume.Matrix {
	f.calls++
	out := dry
	for i := 0; i < out.Channels; i++ {
		out.StaticMatrix[i] *= 0.5
	}
	return out
}

type stubRenderSampler struct {
	done         bool
	active       bool
	muted        bool
	pos          sampling.Pos
	setPosValue  sampling.Pos
	sampleRate   frequency.Frequency
	playbackRate frequency.Frequency
	pan          panning.Position
	period       stubPeriod
	tickCount    int
	sampleCalls  []sampling.Pos
}

func (s *stubRenderSampler) Clone(_ bool) Voice                      { return s }
func (s *stubRenderSampler) DumpState(index.Channel, tracing.Tracer) {}
func (s *stubRenderSampler) Reset() error                            { return nil }
func (s *stubRenderSampler) SetPlaybackRate(rate frequency.Frequency) error {
	s.playbackRate = rate
	return nil
}
func (s *stubRenderSampler) Attack()                            {}
func (s *stubRenderSampler) Release()                           {}
func (s *stubRenderSampler) Fadeout()                           {}
func (s *stubRenderSampler) Stop()                              {}
func (s *stubRenderSampler) Tick() error                        { s.tickCount++; return nil }
func (s *stubRenderSampler) RowEnd() error                      { return nil }
func (s *stubRenderSampler) IsDone() bool                       { return s.done }
func (s *stubRenderSampler) SetMuted(muted bool) error          { s.muted = muted; return nil }
func (s *stubRenderSampler) IsMuted() bool                      { return s.muted }
func (s *stubRenderSampler) GetSampleRate() frequency.Frequency { return s.sampleRate }
func (s *stubRenderSampler) IsActive() bool                     { return s.active }
func (s *stubRenderSampler) SetPos(pos sampling.Pos) error {
	s.setPosValue = pos
	s.pos = pos
	return nil
}
func (s *stubRenderSampler) GetPos() (sampling.Pos, error)       { return s.pos, nil }
func (s *stubRenderSampler) GetFinalPeriod() (stubPeriod, error) { return s.period, nil }
func (s *stubRenderSampler) GetFinalVolume() volume.Volume       { return 1 }
func (s *stubRenderSampler) GetFinalPan() panning.Position       { return s.pan }
func (s *stubRenderSampler) GetSample(pos sampling.Pos) volume.Matrix {
	s.sampleCalls = append(s.sampleCalls, pos)
	v := float32(pos.Pos + 1)
	return volume.Matrix{StaticMatrix: volume.StaticMatrix{volume.Volume(v), volume.Volume(-v)}, Channels: 2}
}

type bareVoice struct{ ticked int }

func (b *bareVoice) Clone(_ bool) Voice                        { return b }
func (b *bareVoice) DumpState(index.Channel, tracing.Tracer)   {}
func (b *bareVoice) Reset() error                              { return nil }
func (b *bareVoice) SetPlaybackRate(frequency.Frequency) error { return nil }
func (b *bareVoice) Attack()                                   {}
func (b *bareVoice) Release()                                  {}
func (b *bareVoice) Fadeout()                                  {}
func (b *bareVoice) Stop()                                     {}
func (b *bareVoice) Tick() error                               { b.ticked++; return nil }
func (b *bareVoice) RowEnd() error                             { return nil }
func (b *bareVoice) IsDone() bool                              { return false }
func (b *bareVoice) SetMuted(bool) error                       { return nil }
func (b *bareVoice) IsMuted() bool                             { return false }
func (b *bareVoice) GetSampleRate() frequency.Frequency        { return 0 }

func TestRenderAndTickReturnsNilWhenDone(t *testing.T) {
	v := &stubRenderSampler{done: true}
	mix := mixing.Mixer{Channels: 2}
	details := mixer.Details{Mix: &mix}

	data, err := RenderAndTick[stubPeriod](v, stubPeriodConverter{samplerAdd: 1}, nil, details, &stubFilter{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data != nil {
		t.Fatalf("expected nil data when voice done")
	}
	if v.tickCount != 0 {
		t.Fatalf("tick should not run when IsDone")
	}
}

func TestRenderAndTickTicksNonRenderSampler(t *testing.T) {
	b := &bareVoice{}

	data, err := RenderAndTick[stubPeriod](b, stubPeriodConverter{samplerAdd: 1}, nil, mixer.Details{}, &stubFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != nil {
		t.Fatalf("expected nil data for non-render sampler")
	}
	if b.ticked != 1 {
		t.Fatalf("tick should be called once, got %d", b.ticked)
	}
}

func TestRenderAndTickRendersAndUpdatesPosition(t *testing.T) {
	centerPan := &stubPanMatrix{factorL: 2, factorR: 3}
	panMixer := &stubDetailPanMixer{matrix: &stubPanMatrix{factorL: 1, factorR: 1}}
	filter := &stubFilter{}

	v := &stubRenderSampler{
		active:     true,
		pos:        sampling.Pos{Pos: 0, Frac: 0},
		sampleRate: 10,
		period:     stubPeriod(4),
		pan:        panning.Position{Angle: 0.1, Distance: 0.5},
	}

	mix := mixing.Mixer{Channels: 2}
	details := mixer.Details{
		Mix:              &mix,
		Panmixer:         panMixer,
		SampleRate:       20,
		StereoSeparation: 0.75,
		Samples:          2,
	}

	data, err := RenderAndTick[stubPeriod](v, stubPeriodConverter{samplerAdd: 1}, centerPan, details, filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data == nil {
		t.Fatalf("expected data output")
	}

	if v.tickCount != 1 {
		t.Fatalf("tick should run once, got %d", v.tickCount)
	}
	if v.playbackRate != details.SampleRate {
		t.Fatalf("playback rate mismatch: got %v want %v", v.playbackRate, details.SampleRate)
	}
	if v.setPosValue != (sampling.Pos{Pos: 2, Frac: 0}) {
		t.Fatalf("expected position to update to Pos 2, got %+v", v.setPosValue)
	}

	if len(v.sampleCalls) != details.Samples {
		t.Fatalf("expected %d samples read, got %d", details.Samples, len(v.sampleCalls))
	}
	if v.sampleCalls[0].Pos != 0 || v.sampleCalls[1].Pos != 1 {
		t.Fatalf("unexpected sample positions: %+v", v.sampleCalls)
	}

	if filter.calls != details.Samples {
		t.Fatalf("filter should be applied per sample, got %d", filter.calls)
	}
	if centerPan.calls != details.Samples {
		t.Fatalf("pan matrix should apply per sample, got %d", centerPan.calls)
	}
	if panMixer.calls != 1 || panMixer.pan != v.pan || panMixer.sep != details.StereoSeparation {
		t.Fatalf("pan mixer should be invoked once with final pan: calls=%d pan=%+v sep=%v", panMixer.calls, panMixer.pan, panMixer.sep)
	}

	if data.PanMatrix != panMixer.matrix {
		t.Fatalf("pan matrix not propagated to output")
	}
	if data.SamplesLen != details.Samples {
		t.Fatalf("sample length mismatch: got %d want %d", data.SamplesLen, details.Samples)
	}
	if data.Volume != 1 {
		t.Fatalf("volume should default to 1, got %v", data.Volume)
	}

	if len(data.Data) != details.Samples {
		t.Fatalf("mixbuffer length mismatch: got %d", len(data.Data))
	}
	first := data.Data[0]
	second := data.Data[1]

	// first sample: dry {1,-1}, filter halves to {0.5,-0.5}, pan scales -> {1,-1.5}
	if first.StaticMatrix[0] != volume.Volume(1) || first.StaticMatrix[1] != volume.Volume(-1.5) {
		t.Fatalf("unexpected first sample mix: %+v", first.StaticMatrix)
	}
	// second sample: dry {2,-2}, filter halves to {1,-1}, pan scales -> {2,-3}
	if second.StaticMatrix[0] != volume.Volume(2) || second.StaticMatrix[1] != volume.Volume(-3) {
		t.Fatalf("unexpected second sample mix: %+v", second.StaticMatrix)
	}
}
