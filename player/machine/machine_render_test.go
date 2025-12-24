package machine

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/mixing"
	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/player/sampler"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/system"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/mixer"
	"github.com/gotracker/playback/voice/vol0optimization"
)

type stubGV float32

func (stubGV) IsInvalid() bool           { return false }
func (stubGV) IsUseInstrumentVol() bool  { return false }
func (g stubGV) ToVolume() volume.Volume { return volume.Volume(g) }

type stubPan float32

func (stubPan) IsInvalid() bool { return false }
func (p stubPan) ToPosition() panning.Position {
	return panning.Position{Angle: 0, Distance: float32(p)}
}

type stubPeriod struct{}

func (stubPeriod) IsInvalid() bool { return false }

type stubPeriodCalc struct{}

func (stubPeriodCalc) GetSystem() system.System         { return nil }
func (stubPeriodCalc) GetPeriod(_ note.Note) stubPeriod { return stubPeriod{} }
func (stubPeriodCalc) PortaToNote(p stubPeriod, _ period.Delta, _ stubPeriod) (stubPeriod, error) {
	return p, nil
}
func (stubPeriodCalc) PortaDown(p stubPeriod, _ period.Delta) (stubPeriod, error) { return p, nil }
func (stubPeriodCalc) PortaUp(p stubPeriod, _ period.Delta) (stubPeriod, error)   { return p, nil }
func (stubPeriodCalc) AddDelta(p stubPeriod, _ period.Delta) (stubPeriod, error)  { return p, nil }
func (stubPeriodCalc) GetSamplerAdd(stubPeriod, frequency.Frequency, frequency.Frequency) float64 {
	return 0
}
func (stubPeriodCalc) GetFrequency(stubPeriod) frequency.Frequency { return 0 }

type stubChannelMemory struct{}

func (stubChannelMemory) Retrigger()   {}
func (stubChannelMemory) StartOrder0() {}

type stubChannelSettings struct{}

func (stubChannelSettings) IsEnabled() bool                   { return true }
func (stubChannelSettings) IsMuted() bool                     { return false }
func (stubChannelSettings) GetOutputChannelNum() int          { return 0 }
func (stubChannelSettings) GetMemory() song.ChannelMemory     { return stubChannelMemory{} }
func (stubChannelSettings) IsPanEnabled() bool                { return true }
func (stubChannelSettings) GetDefaultFilterInfo() filter.Info { return filter.Info{} }
func (stubChannelSettings) IsDefaultFilterEnabled() bool      { return false }
func (stubChannelSettings) GetVol0OptimizationSettings() vol0optimization.Vol0OptimizationSettings {
	return vol0optimization.Vol0OptimizationSettings{}
}
func (stubChannelSettings) GetOPLChannel() index.OPLChannel { return 0 }

type stubSongData struct{}

func (stubSongData) GetPeriodType() reflect.Type              { return reflect.TypeOf(stubPeriod{}) }
func (stubSongData) GetGlobalVolumeType() reflect.Type        { return reflect.TypeOf(stubGV(0)) }
func (stubSongData) GetChannelMixingVolumeType() reflect.Type { return reflect.TypeOf(stubGV(0)) }
func (stubSongData) GetChannelVolumeType() reflect.Type       { return reflect.TypeOf(stubGV(0)) }
func (stubSongData) GetChannelPanningType() reflect.Type      { return reflect.TypeOf(stubPan(0)) }
func (stubSongData) GetInitialBPM() int                       { return 6 }
func (stubSongData) GetInitialTempo() int                     { return 6 }
func (stubSongData) GetMixingVolumeGeneric() volume.Volume    { return 1 }
func (stubSongData) GetTickDuration(int) time.Duration        { return time.Second }
func (stubSongData) GetOrderList() []index.Pattern            { return []index.Pattern{0} }
func (stubSongData) GetNumChannels() int                      { return 1 }
func (stubSongData) GetChannelSettings(index.Channel) song.ChannelSettings {
	return stubChannelSettings{}
}
func (stubSongData) NumInstruments() int { return 0 }
func (stubSongData) GetInstrument(int, note.Semitone) (instrument.InstrumentIntf, note.Semitone) {
	return nil, 0
}
func (stubSongData) GetName() string { return "stub" }
func (stubSongData) GetPatternByOrder(index.Order) (song.Pattern, error) {
	return nil, song.ErrStopSong
}
func (stubSongData) GetPattern(index.Pattern) (song.Pattern, error)            { return nil, song.ErrStopSong }
func (stubSongData) GetPeriodCalculator() song.PeriodCalculatorIntf            { return nil }
func (stubSongData) GetInitialOrder() index.Order                              { return 0 }
func (stubSongData) GetRowRenderStringer(song.Row, int, bool) song.RowStringer { return nil }
func (stubSongData) GetSystem() system.System                                  { return nil }
func (stubSongData) GetMachineSettings() any                                   { return nil }
func (stubSongData) ForEachChannel(bool, func(index.Channel) (bool, error)) error {
	return nil
}
func (stubSongData) IsOPL2Enabled() bool { return false }

type zeroTickSongData struct{ stubSongData }

func (zeroTickSongData) GetTickDuration(int) time.Duration { return 0 }

type tick1500SongData struct{ stubSongData }

func (tick1500SongData) GetTickDuration(int) time.Duration { return 1500 * time.Millisecond }

type doneVoice struct{ ticked int }

func (v *doneVoice) Clone(bool) voice.Voice                    { return v }
func (v *doneVoice) DumpState(index.Channel, tracing.Tracer)   {}
func (v *doneVoice) Reset() error                              { return nil }
func (v *doneVoice) SetPlaybackRate(frequency.Frequency) error { return nil }
func (v *doneVoice) Attack()                                   {}
func (v *doneVoice) Release()                                  {}
func (v *doneVoice) Fadeout()                                  {}
func (v *doneVoice) Stop()                                     {}
func (v *doneVoice) Tick() error                               { v.ticked++; return nil }
func (v *doneVoice) RowEnd() error                             { return nil }
func (v *doneVoice) IsDone() bool                              { return true }
func (v *doneVoice) SetMuted(bool) error                       { return nil }
func (v *doneVoice) IsMuted() bool                             { return false }
func (v *doneVoice) GetSampleRate() frequency.Frequency        { return 0 }

type errorRenderVoice struct {
	pos sampling.Pos
	err error
}

func (v *errorRenderVoice) Clone(bool) voice.Voice                    { return v }
func (v *errorRenderVoice) DumpState(index.Channel, tracing.Tracer)   {}
func (v *errorRenderVoice) Reset() error                              { return nil }
func (v *errorRenderVoice) SetPlaybackRate(frequency.Frequency) error { return nil }
func (v *errorRenderVoice) Attack()                                   {}
func (v *errorRenderVoice) Release()                                  {}
func (v *errorRenderVoice) Fadeout()                                  {}
func (v *errorRenderVoice) Stop()                                     {}
func (v *errorRenderVoice) Tick() error                               { return nil }
func (v *errorRenderVoice) RowEnd() error                             { return nil }
func (v *errorRenderVoice) IsDone() bool                              { return false }
func (v *errorRenderVoice) SetMuted(bool) error                       { return nil }
func (v *errorRenderVoice) IsMuted() bool                             { return false }
func (v *errorRenderVoice) GetSampleRate() frequency.Frequency        { return 1 }
func (v *errorRenderVoice) GetSample(sampling.Pos) volume.Matrix      { return volume.Matrix{} }
func (v *errorRenderVoice) IsActive() bool                            { return true }
func (v *errorRenderVoice) SetPos(pos sampling.Pos) error             { v.pos = pos; return nil }
func (v *errorRenderVoice) GetPos() (sampling.Pos, error)             { return v.pos, nil }
func (v *errorRenderVoice) GetFinalPeriod() (stubPeriod, error)       { return stubPeriod{}, v.err }
func (v *errorRenderVoice) GetFinalVolume() volume.Volume             { return 1 }
func (v *errorRenderVoice) GetFinalPan() panning.Position             { return panning.CenterAhead }

type stubHardwareSynth struct {
	called  int
	center  panning.PanMixer
	details mixer.Details
}

func (s *stubHardwareSynth) RenderTick(centerAheadPan panning.PanMixer, details mixer.Details) (mixing.Data, mixerVolumeAdjuster, error) {
	s.called++
	s.center = centerAheadPan
	s.details = details

	data := mixing.Data{
		Data:       details.Mix.NewMixBuffer(details.Samples),
		PanMatrix:  centerAheadPan,
		Volume:     volume.Volume(0.25),
		SamplesLen: details.Samples,
	}

	return data, func(v volume.Volume) volume.Volume {
		return v * 0.5
	}, nil
}

type errorHardwareSynth struct{ err error }

func (s errorHardwareSynth) RenderTick(centerAheadPan panning.PanMixer, details mixer.Details) (mixing.Data, mixerVolumeAdjuster, error) {
	return mixing.Data{}, nil, s.err
}

type nilAdjustHardwareSynth struct {
	called bool
}

func (s *nilAdjustHardwareSynth) RenderTick(centerAheadPan panning.PanMixer, details mixer.Details) (mixing.Data, mixerVolumeAdjuster, error) {
	s.called = true
	return mixing.Data{
		Data:       details.Mix.NewMixBuffer(details.Samples),
		PanMatrix:  centerAheadPan,
		Volume:     volume.Volume(0.75),
		SamplesLen: details.Samples,
	}, nil, nil
}

type sizedHardwareSynth struct {
	samples int
}

func (s sizedHardwareSynth) RenderTick(centerAheadPan panning.PanMixer, details mixer.Details) (mixing.Data, mixerVolumeAdjuster, error) {
	return mixing.Data{
		Data:       details.Mix.NewMixBuffer(s.samples),
		PanMatrix:  centerAheadPan,
		Volume:     volume.Volume(0.9),
		SamplesLen: s.samples,
	}, nil, nil
}

func TestPrepareRenderFrameCopiesRowStringerOnTickZero(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		globals: globals[stubGV]{
			bpm:   6,
			tempo: 6,
			gv:    stubGV(1),
			mv:    1,
		},
		ms:          &settings.MachineSettings[stubPeriod, stubGV, stubGV, stubGV, stubPan]{PeriodConverter: stubPeriodCalc{}},
		songData:    stubSongData{},
		rowStringer: stubRowStringer{s: "row text"},
	}
	m.ticker.current = Position{Order: 2, Row: 3, Tick: 0}

	s := sampler.NewSampler(10, 2, 1, nil)

	frame, err := m.prepareRenderFrame(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frame.renderRow.RowText != m.rowStringer {
		t.Fatalf("expected row stringer to be copied on tick 0")
	}
	userdata, ok := frame.premix.Userdata.(*render.RowRender)
	if !ok {
		t.Fatalf("expected premix userdata to store row render, got %T", frame.premix.Userdata)
	}
	if userdata.RowText != m.rowStringer {
		t.Fatalf("expected userdata row stringer to be copied")
	}
	if frame.renderRow.Order != int(m.ticker.current.Order) || frame.renderRow.Row != int(m.ticker.current.Row) || frame.renderRow.Tick != m.ticker.current.Tick {
		t.Fatalf("row render metadata mismatch: got %+v", frame.renderRow)
	}
}

func TestPrepareRenderFrameSkipsRowStringerOnNonZeroTick(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		globals: globals[stubGV]{
			bpm:   6,
			tempo: 6,
			gv:    stubGV(1),
			mv:    1,
		},
		ms:          &settings.MachineSettings[stubPeriod, stubGV, stubGV, stubGV, stubPan]{PeriodConverter: stubPeriodCalc{}},
		songData:    stubSongData{},
		rowStringer: stubRowStringer{s: "row text"},
	}
	m.ticker.current = Position{Order: 1, Row: 4, Tick: 2}

	s := sampler.NewSampler(10, 2, 1, nil)

	frame, err := m.prepareRenderFrame(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frame.renderRow.RowText != nil {
		t.Fatalf("expected row stringer to be nil on non-zero tick")
	}
	userdata, ok := frame.premix.Userdata.(*render.RowRender)
	if !ok {
		t.Fatalf("expected premix userdata to store row render, got %T", frame.premix.Userdata)
	}
	if userdata.RowText != nil {
		t.Fatalf("expected userdata row stringer to be nil on non-zero tick")
	}
}

func TestPrepareRenderFrameUsesTickDurationForSamples(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		globals: globals[stubGV]{
			bpm:   6,
			tempo: 6,
			gv:    stubGV(1),
			mv:    1,
		},
		ms:       &settings.MachineSettings[stubPeriod, stubGV, stubGV, stubGV, stubPan]{PeriodConverter: stubPeriodCalc{}},
		songData: tick1500SongData{},
	}

	s := sampler.NewSampler(10, 2, 1, nil)

	frame, err := m.prepareRenderFrame(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frame.premix.SamplesLen != 15 {
		t.Fatalf("expected samples len 15 from 1.5s tick at 10Hz, got %d", frame.premix.SamplesLen)
	}
	if frame.details.Samples != 15 {
		t.Fatalf("expected details samples 15, got %d", frame.details.Samples)
	}
	if frame.details.Duration != 1500*time.Millisecond {
		t.Fatalf("unexpected duration: %v", frame.details.Duration)
	}
}

func TestNormalizePremixAddsSilence(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{}

	panMixer := mixing.GetPanMixer(2)
	if panMixer == nil {
		t.Fatalf("expected stereo pan mixer")
	}

	details := mixer.Details{
		Mix:              &mixing.Mixer{Channels: 2},
		Panmixer:         panMixer,
		SampleRate:       10,
		StereoSeparation: 1,
		Samples:          3,
	}

	center := panMixer.GetMixingMatrix(panning.CenterAhead, details.StereoSeparation)

	frame := renderFrame{
		details:        details,
		centerAheadPan: center,
		premix:         output.PremixData{},
	}

	m.normalizePremix(&frame)

	if len(frame.premix.Data) != 1 {
		t.Fatalf("expected premix to contain one channel, got %d", len(frame.premix.Data))
	}
	if len(frame.premix.Data[0]) != 1 {
		t.Fatalf("expected channel data to have one entry, got %d", len(frame.premix.Data[0]))
	}
	entry := frame.premix.Data[0][0]
	if entry.SamplesLen != details.Samples {
		t.Fatalf("unexpected samples len: got %d want %d", entry.SamplesLen, details.Samples)
	}
	if len(entry.Data) != details.Samples {
		t.Fatalf("mixbuffer length mismatch: got %d", len(entry.Data))
	}
}

func TestNormalizePremixZeroSamples(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{}

	panMixer := mixing.GetPanMixer(2)
	if panMixer == nil {
		t.Fatalf("expected stereo pan mixer")
	}

	details := mixer.Details{
		Mix:              &mixing.Mixer{Channels: 2},
		Panmixer:         panMixer,
		SampleRate:       10,
		StereoSeparation: 1,
		Samples:          0,
	}

	frame := renderFrame{
		details:        details,
		centerAheadPan: panMixer.GetMixingMatrix(panning.CenterAhead, details.StereoSeparation),
		premix:         output.PremixData{},
	}

	m.normalizePremix(&frame)

	if len(frame.premix.Data) != 1 || len(frame.premix.Data[0]) != 1 {
		t.Fatalf("expected single channel entry even for zero samples")
	}
	entry := frame.premix.Data[0][0]
	if entry.SamplesLen != 0 {
		t.Fatalf("expected samples len 0, got %d", entry.SamplesLen)
	}
	if len(entry.Data) != 0 {
		t.Fatalf("expected zero-length mixbuffer, got %d", len(entry.Data))
	}
}

func TestNormalizePremixKeepsExistingData(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{}

	panMixer := mixing.GetPanMixer(2)
	if panMixer == nil {
		t.Fatalf("expected stereo pan mixer")
	}

	details := mixer.Details{
		Mix:              &mixing.Mixer{Channels: 2},
		Panmixer:         panMixer,
		SampleRate:       10,
		StereoSeparation: 1,
		Samples:          2,
	}

	existing := mixing.ChannelData{mixing.Data{SamplesLen: details.Samples}}
	frame := renderFrame{
		details:        details,
		centerAheadPan: panMixer.GetMixingMatrix(panning.CenterAhead, details.StereoSeparation),
		premix: output.PremixData{
			SamplesLen: details.Samples,
			Data:       []mixing.ChannelData{existing},
		},
	}

	m.normalizePremix(&frame)

	if len(frame.premix.Data) != 1 {
		t.Fatalf("expected existing data retained")
	}
	if frame.premix.Data[0][0].SamplesLen != details.Samples {
		t.Fatalf("unexpected samples len after normalize: %d", frame.premix.Data[0][0].SamplesLen)
	}
}

func TestRenderVoicesGeneratesSilence(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		ms: &settings.MachineSettings[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
			PeriodConverter: stubPeriodCalc{},
		},
	}

	panMixer := mixing.GetPanMixer(2)
	if panMixer == nil {
		t.Fatalf("expected stereo pan mixer")
	}

	details := mixer.Details{
		Mix:              &mixing.Mixer{Channels: 2},
		Panmixer:         panMixer,
		SampleRate:       10,
		StereoSeparation: 1,
		Samples:          2,
	}
	center := panMixer.GetMixingMatrix(panning.CenterAhead, details.StereoSeparation)

	ch := render.Channel[stubPeriod]{}
	ch.StartVoice(&doneVoice{}, nil)
	m.actualOutputs = []render.Channel[stubPeriod]{ch}
	m.virtualOutputs = []render.Channel[stubPeriod]{{}}

	frame := renderFrame{
		details:        details,
		centerAheadPan: center,
	}

	mixData, err := m.renderVoices(frame)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := len(m.actualOutputs) + len(m.virtualOutputs)
	if len(mixData) != expected {
		t.Fatalf("expected %d mix entries, got %d", expected, len(mixData))
	}

	for i, d := range mixData {
		if d.Volume != 0 {
			t.Fatalf("entry %d expected volume 0, got %v", i, d.Volume)
		}
		if applied := d.PanMatrix.Apply(volume.Volume(1)); applied.Channels != 2 {
			t.Fatalf("entry %d expected stereo pan matrix, got %d channels", i, applied.Channels)
		}
		if d.SamplesLen != details.Samples {
			t.Fatalf("entry %d samples len mismatch: got %d want %d", i, d.SamplesLen, details.Samples)
		}
		if len(d.Data) != details.Samples {
			t.Fatalf("entry %d mixbuffer length mismatch: got %d", i, len(d.Data))
		}
	}
}

func TestRenderVoicesZeroSamples(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		ms: &settings.MachineSettings[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
			PeriodConverter: stubPeriodCalc{},
		},
	}

	panMixer := mixing.GetPanMixer(2)
	if panMixer == nil {
		t.Fatalf("expected stereo pan mixer")
	}

	details := mixer.Details{
		Mix:              &mixing.Mixer{Channels: 2},
		Panmixer:         panMixer,
		SampleRate:       10,
		StereoSeparation: 1,
		Samples:          0,
	}
	center := panMixer.GetMixingMatrix(panning.CenterAhead, details.StereoSeparation)

	ch := render.Channel[stubPeriod]{}
	ch.StartVoice(&doneVoice{}, nil)
	m.actualOutputs = []render.Channel[stubPeriod]{ch}
	m.virtualOutputs = []render.Channel[stubPeriod]{{}}

	frame := renderFrame{
		details:        details,
		centerAheadPan: center,
	}

	mixData, err := m.renderVoices(frame)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, d := range mixData {
		if d.SamplesLen != 0 {
			t.Fatalf("entry %d expected samples len 0, got %d", i, d.SamplesLen)
		}
		if len(d.Data) != 0 {
			t.Fatalf("entry %d expected zero-length mixbuffer, got %d", i, len(d.Data))
		}
	}
}

func TestRenderVoicesPropagatesVoiceError(t *testing.T) {
	errRender := errors.New("render fail")
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		ms: &settings.MachineSettings[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
			PeriodConverter: stubPeriodCalc{},
		},
	}

	panMixer := mixing.GetPanMixer(2)
	if panMixer == nil {
		t.Fatalf("expected stereo pan mixer")
	}

	details := mixer.Details{
		Mix:              &mixing.Mixer{Channels: 2},
		Panmixer:         panMixer,
		SampleRate:       10,
		StereoSeparation: 1,
		Samples:          2,
	}

	ch := render.Channel[stubPeriod]{}
	ch.StartVoice(&errorRenderVoice{err: errRender}, nil)
	m.actualOutputs = []render.Channel[stubPeriod]{ch}

	frame := renderFrame{
		details:        details,
		centerAheadPan: panMixer.GetMixingMatrix(panning.CenterAhead, details.StereoSeparation),
	}

	if _, err := m.renderVoices(frame); !errors.Is(err, errRender) {
		t.Fatalf("expected render error to propagate, got %v", err)
	}
}

func TestRenderInvokesOnGenerateWithPremix(t *testing.T) {
	called := 0
	var gotPremix *output.PremixData
	s := sampler.NewSampler(10, 2, 1, func(p *output.PremixData) {
		called++
		gotPremix = p
	})

	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		globals: globals[stubGV]{
			bpm:   6,
			tempo: 6,
			gv:    stubGV(1),
			mv:    0.5,
		},
		ms:       &settings.MachineSettings[stubPeriod, stubGV, stubGV, stubGV, stubPan]{PeriodConverter: stubPeriodCalc{}},
		us:       settings.UserSettings{},
		songData: stubSongData{},
	}

	if err := m.Render(s); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if called != 1 {
		t.Fatalf("expected OnGenerate to be called once, got %d", called)
	}
	if gotPremix == nil {
		t.Fatalf("expected premix data")
	}

	expectedLen := int(float64(s.SampleRate) * float64(time.Second) / float64(time.Second))
	if gotPremix.SamplesLen != expectedLen {
		t.Fatalf("unexpected samples len: got %d want %d", gotPremix.SamplesLen, expectedLen)
	}

	expectedVol := m.gv.ToVolume() * m.mv
	if gotPremix.MixerVolume != expectedVol {
		t.Fatalf("unexpected mixer volume: got %v want %v", gotPremix.MixerVolume, expectedVol)
	}

	if len(gotPremix.Data) != 1 {
		t.Fatalf("expected premix to contain one channel, got %d", len(gotPremix.Data))
	}
	if len(gotPremix.Data[0]) != 1 {
		t.Fatalf("expected channel data to have one entry, got %d", len(gotPremix.Data[0]))
	}
}

func TestRenderFailsOnInvalidTickDuration(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		globals: globals[stubGV]{
			bpm:   6,
			tempo: 6,
			gv:    stubGV(1),
			mv:    0.5,
		},
		ms:       &settings.MachineSettings[stubPeriod, stubGV, stubGV, stubGV, stubPan]{PeriodConverter: stubPeriodCalc{}},
		songData: zeroTickSongData{},
	}

	s := sampler.NewSampler(10, 2, 1, nil)

	if err := m.Render(s); err == nil {
		t.Fatalf("expected error when tick duration is invalid")
	}
}

func TestRenderPropagatesRenderVoicesError(t *testing.T) {
	errRender := errors.New("render voices fail")

	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		globals: globals[stubGV]{
			bpm:   6,
			tempo: 6,
			gv:    stubGV(1),
			mv:    0.5,
		},
		ms: &settings.MachineSettings[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
			PeriodConverter: stubPeriodCalc{},
		},
		songData: stubSongData{},
	}

	ch := render.Channel[stubPeriod]{}
	ch.StartVoice(&errorRenderVoice{err: errRender}, nil)
	m.actualOutputs = []render.Channel[stubPeriod]{ch}

	s := sampler.NewSampler(10, 2, 1, nil)

	if err := m.Render(s); !errors.Is(err, errRender) {
		t.Fatalf("expected renderVoices error to propagate, got %v", err)
	}
}

func TestMixHardwareSynthsAppendsDataAndAdjustsVolume(t *testing.T) {
	panMixer := mixing.GetPanMixer(2)
	if panMixer == nil {
		t.Fatalf("expected stereo pan mixer")
	}

	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{}
	m.hardwareSynths = []hardwareSynth{&stubHardwareSynth{}}

	premix := output.PremixData{
		SamplesLen:  4,
		MixerVolume: 1,
	}

	details := mixer.Details{
		Mix:              &mixing.Mixer{Channels: 2},
		Panmixer:         panMixer,
		SampleRate:       10,
		StereoSeparation: 1,
		Samples:          premix.SamplesLen,
	}

	frame := renderFrame{
		details:        details,
		centerAheadPan: panMixer.GetMixingMatrix(panning.CenterAhead, details.StereoSeparation),
		premix:         premix,
	}

	if err := m.mixHardwareSynths(frame, &frame.premix); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frame.premix.MixerVolume != 0.5 {
		t.Fatalf("expected mixer volume adjusted to 0.5, got %v", frame.premix.MixerVolume)
	}

	if len(frame.premix.Data) != 1 {
		t.Fatalf("expected 1 channel data entry, got %d", len(frame.premix.Data))
	}
	channel := frame.premix.Data[0]
	if len(channel) != 1 {
		t.Fatalf("expected 1 mix data entry, got %d", len(channel))
	}
	d := channel[0]
	if d.SamplesLen != premix.SamplesLen {
		t.Fatalf("samples len mismatch: got %d want %d", d.SamplesLen, premix.SamplesLen)
	}
	if d.Volume != volume.Volume(0.25) {
		t.Fatalf("unexpected volume: got %v want %v", d.Volume, volume.Volume(0.25))
	}
	if len(d.Data) != premix.SamplesLen {
		t.Fatalf("mixbuffer length mismatch: got %d want %d", len(d.Data), premix.SamplesLen)
	}
	if applied := d.PanMatrix.Apply(volume.Volume(1)); applied.Channels != 2 {
		t.Fatalf("expected stereo pan matrix, got %d channels", applied.Channels)
	}
}

func TestMixHardwareSynthsNoDevicesLeavesPremix(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{}

	premix := output.PremixData{SamplesLen: 2, MixerVolume: 0.6}

	details := mixer.Details{
		Mix:              &mixing.Mixer{Channels: 2},
		Panmixer:         mixing.GetPanMixer(2),
		SampleRate:       10,
		StereoSeparation: 1,
		Samples:          premix.SamplesLen,
	}

	frame := renderFrame{details: details, premix: premix}

	if err := m.mixHardwareSynths(frame, &frame.premix); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if frame.premix.MixerVolume != premix.MixerVolume {
		t.Fatalf("expected mixer volume unchanged, got %v", frame.premix.MixerVolume)
	}
	if len(frame.premix.Data) != 0 {
		t.Fatalf("expected no data appended, got %d entries", len(frame.premix.Data))
	}
}

func TestMixHardwareSynthsReturnsError(t *testing.T) {
	panMixer := mixing.GetPanMixer(2)
	if panMixer == nil {
		t.Fatalf("expected stereo pan mixer")
	}

	errRender := errors.New("render failure")
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{}
	m.hardwareSynths = []hardwareSynth{errorHardwareSynth{err: errRender}}

	premix := output.PremixData{
		SamplesLen:  2,
		MixerVolume: 1,
	}

	details := mixer.Details{
		Mix:              &mixing.Mixer{Channels: 2},
		Panmixer:         panMixer,
		SampleRate:       10,
		StereoSeparation: 1,
		Samples:          premix.SamplesLen,
	}

	frame := renderFrame{
		details:        details,
		centerAheadPan: panMixer.GetMixingMatrix(panning.CenterAhead, details.StereoSeparation),
		premix:         premix,
	}

	err := m.mixHardwareSynths(frame, &premix)
	if !errors.Is(err, errRender) {
		t.Fatalf("expected error to propagate, got %v", err)
	}
	if len(premix.Data) != 0 {
		t.Fatalf("expected premix data to remain unchanged on error, got %d entries", len(premix.Data))
	}
	if premix.MixerVolume != 1 {
		t.Fatalf("expected mixer volume unchanged, got %v", premix.MixerVolume)
	}
}

func TestMixHardwareSynthsKeepsVolumeWhenNoAdjuster(t *testing.T) {
	panMixer := mixing.GetPanMixer(2)
	if panMixer == nil {
		t.Fatalf("expected stereo pan mixer")
	}

	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{}
	var synth nilAdjustHardwareSynth
	m.hardwareSynths = []hardwareSynth{&synth}

	premix := output.PremixData{
		SamplesLen:  3,
		MixerVolume: 0.8,
	}

	details := mixer.Details{
		Mix:              &mixing.Mixer{Channels: 2},
		Panmixer:         panMixer,
		SampleRate:       10,
		StereoSeparation: 1,
		Samples:          premix.SamplesLen,
	}

	frame := renderFrame{
		details:        details,
		centerAheadPan: panMixer.GetMixingMatrix(panning.CenterAhead, details.StereoSeparation),
		premix:         premix,
	}

	if err := m.mixHardwareSynths(frame, &premix); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !synth.called {
		t.Fatalf("expected hardware synth to be called")
	}
	if premix.MixerVolume != 0.8 {
		t.Fatalf("expected mixer volume to remain unchanged, got %v", premix.MixerVolume)
	}
	if len(premix.Data) != 1 {
		t.Fatalf("expected 1 channel data entry, got %d", len(premix.Data))
	}
	channel := premix.Data[0]
	if len(channel) != 1 {
		t.Fatalf("expected 1 mix data entry, got %d", len(channel))
	}
	if channel[0].Volume != volume.Volume(0.75) {
		t.Fatalf("unexpected volume: got %v want %v", channel[0].Volume, volume.Volume(0.75))
	}
}

func TestMixHardwareSynthsAppliesAllAdjusters(t *testing.T) {
	panMixer := mixing.GetPanMixer(2)
	if panMixer == nil {
		t.Fatalf("expected stereo pan mixer")
	}

	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{}
	m.hardwareSynths = []hardwareSynth{&stubHardwareSynth{}, &stubHardwareSynth{}}

	premix := output.PremixData{SamplesLen: 2, MixerVolume: 1}

	details := mixer.Details{
		Mix:              &mixing.Mixer{Channels: 2},
		Panmixer:         panMixer,
		SampleRate:       10,
		StereoSeparation: 1,
		Samples:          premix.SamplesLen,
	}

	frame := renderFrame{
		details:        details,
		centerAheadPan: panMixer.GetMixingMatrix(panning.CenterAhead, details.StereoSeparation),
		premix:         premix,
	}

	if err := m.mixHardwareSynths(frame, &frame.premix); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if frame.premix.MixerVolume != 0.25 {
		t.Fatalf("expected mixer volume adjusted twice to 0.25, got %v", frame.premix.MixerVolume)
	}
	if len(frame.premix.Data) != 2 {
		t.Fatalf("expected 2 channel data entries, got %d", len(frame.premix.Data))
	}
}

func TestMixHardwareSynthsAllowsMismatchedSampleLens(t *testing.T) {
	panMixer := mixing.GetPanMixer(2)
	if panMixer == nil {
		t.Fatalf("expected stereo pan mixer")
	}

	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{}
	m.hardwareSynths = []hardwareSynth{sizedHardwareSynth{samples: 3}}

	premix := output.PremixData{SamplesLen: 2, MixerVolume: 0.7}

	details := mixer.Details{
		Mix:              &mixing.Mixer{Channels: 2},
		Panmixer:         panMixer,
		SampleRate:       10,
		StereoSeparation: 1,
		Samples:          premix.SamplesLen,
	}

	frame := renderFrame{
		details:        details,
		centerAheadPan: panMixer.GetMixingMatrix(panning.CenterAhead, details.StereoSeparation),
		premix:         premix,
	}

	if err := m.mixHardwareSynths(frame, &frame.premix); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if frame.premix.MixerVolume != premix.MixerVolume {
		t.Fatalf("expected mixer volume unchanged, got %v", frame.premix.MixerVolume)
	}
	if len(frame.premix.Data) != 1 {
		t.Fatalf("expected 1 channel data entry, got %d", len(frame.premix.Data))
	}
	if got := frame.premix.Data[0][0].SamplesLen; got != 3 {
		t.Fatalf("expected appended data to keep synth samples len 3, got %d", got)
	}
}
