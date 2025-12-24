package machine

import (
	"errors"
	"testing"
	"time"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/player/sampler"
	"github.com/gotracker/playback/song"
	optional "github.com/heucuva/optional"
)

func TestTickRequiresSampler(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{}

	if err := m.Tick(nil); err == nil {
		t.Fatalf("expected error when sampler is nil")
	}
}

func TestRenderRequiresSampler(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{}

	if err := m.Render(nil); err == nil {
		t.Fatalf("expected error when sampler is nil")
	}
}

func TestAdvanceStopsWhenSongStops(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		songData: stubSongData{},
	}

	err := m.Advance()
	if !errors.Is(err, song.ErrStopSong) {
		t.Fatalf("expected song.ErrStopSong, got %v", err)
	}
}

func TestExtraTicksAndRepeats(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{}

	if err := m.AddExtraTicks(2); err != nil {
		t.Fatalf("AddExtraTicks unexpected error: %v", err)
	}
	if m.extraTicks != 2 {
		t.Fatalf("expected extraTicks=2, got %d", m.extraTicks)
	}

	if err := m.RowRepeat(3); err != nil {
		t.Fatalf("RowRepeat unexpected error: %v", err)
	}
	if m.repeats != 3 {
		t.Fatalf("expected repeats=3, got %d", m.repeats)
	}

	if !m.consumeRepeat() {
		t.Fatalf("expected consumeRepeat to report true")
	}
	if m.repeats != 2 {
		t.Fatalf("expected repeats decremented to 2, got %d", m.repeats)
	}

	if err := m.AddExtraTicks(-1); err == nil {
		t.Fatalf("expected error for negative extra ticks")
	}
	if err := m.RowRepeat(-1); err == nil {
		t.Fatalf("expected error for negative repeats")
	}
}

type seqSongData struct {
	stubSongData
	pat song.Pattern
}

func (s seqSongData) GetPatternByOrder(index.Order) (song.Pattern, error) { return s.pat, nil }
func (s seqSongData) GetPattern(index.Pattern) (song.Pattern, error)      { return s.pat, nil }
func (s seqSongData) GetOrderList() []index.Pattern                       { return []index.Pattern{0} }
func (s seqSongData) ForEachChannel(bool, func(index.Channel) (bool, error)) error {
	return nil
}
func (s seqSongData) GetRowRenderStringer(song.Row, int, bool) song.RowStringer {
	return stubRowStringer{s: "row"}
}

type stopSongData struct {
	stubSongData
}

func (stopSongData) ForEachChannel(bool, func(index.Channel) (bool, error)) error {
	return song.ErrStopSong
}

type playSongData struct {
	stubSongData
	pat song.Pattern
}

func (p playSongData) GetPatternByOrder(index.Order) (song.Pattern, error) { return p.pat, nil }
func (p playSongData) GetPattern(index.Pattern) (song.Pattern, error)      { return p.pat, nil }
func (p playSongData) GetOrderList() []index.Pattern                       { return []index.Pattern{0} }
func (p playSongData) ForEachChannel(bool, func(index.Channel) (bool, error)) error {
	return nil
}
func (p playSongData) GetRowRenderStringer(song.Row, int, bool) song.RowStringer {
	return stubRowStringer{s: "row"}
}
func (p playSongData) GetTickDuration(int) time.Duration { return time.Second }

func TestAdvanceProgressesRowsAndLoopsOrder(t *testing.T) {
	pat := song.Pattern{song.Row(0), song.Row(1)}
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		ticker:   ticker{},
		globals:  globals[stubGV]{tempo: 2},
		songData: seqSongData{pat: pat},
	}

	if err := initTick(&m.ticker, &m, tickerSettings{InitialOrder: 0, InitialRow: 0, SongLoopCount: -1}); err != nil {
		t.Fatalf("initTick error: %v", err)
	}

	if err := m.Advance(); err != nil {
		t.Fatalf("advance1 error: %v", err)
	}
	if m.ticker.current.Row != 0 || m.ticker.current.Tick != 1 || m.ticker.current.Order != 0 {
		t.Fatalf("state after advance1 got row %d tick %d order %d", m.ticker.current.Row, m.ticker.current.Tick, m.ticker.current.Order)
	}

	if err := m.Advance(); err != nil {
		t.Fatalf("advance2 error: %v", err)
	}
	if m.ticker.current.Row != 1 || m.ticker.current.Tick != 0 || m.ticker.current.Order != 0 {
		t.Fatalf("state after advance2 got row %d tick %d order %d", m.ticker.current.Row, m.ticker.current.Tick, m.ticker.current.Order)
	}

	if err := m.Advance(); err != nil {
		t.Fatalf("advance3 error: %v", err)
	}
	if m.ticker.current.Row != 1 || m.ticker.current.Tick != 1 || m.ticker.current.Order != 0 {
		t.Fatalf("state after advance3 got row %d tick %d order %d", m.ticker.current.Row, m.ticker.current.Tick, m.ticker.current.Order)
	}

	if err := m.Advance(); err != nil {
		t.Fatalf("advance4 error: %v", err)
	}
	if m.ticker.current.Row != 0 || m.ticker.current.Tick != 0 || m.ticker.current.Order != 0 {
		t.Fatalf("state after advance4 got row %d tick %d order %d", m.ticker.current.Row, m.ticker.current.Tick, m.ticker.current.Order)
	}
}

func TestAdvanceStopsOnSongLoopCountZero(t *testing.T) {
	pat := song.Pattern{song.Row(0)}
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		ticker:   ticker{},
		globals:  globals[stubGV]{tempo: 1},
		songData: seqSongData{pat: pat},
	}

	if err := initTick(&m.ticker, &m, tickerSettings{InitialOrder: 0, InitialRow: 0, SongLoopCount: 0, SongLoopStartingOrder: 0}); err != nil {
		t.Fatalf("initTick error: %v", err)
	}

	err := m.Advance()
	if !errors.Is(err, song.ErrStopSong) {
		t.Fatalf("expected ErrStopSong due to loop count, got %v", err)
	}
}

func TestAdvanceStopsAtPlayUntilPosition(t *testing.T) {
	pat := song.Pattern{song.Row(0), song.Row(1)}
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		ticker:   ticker{},
		globals:  globals[stubGV]{tempo: 1},
		songData: seqSongData{pat: pat},
	}

	playUntilOrder := optional.NewValue[index.Order](0)
	playUntilRow := optional.NewValue[index.Row](1)

	ts := tickerSettings{InitialOrder: 0, InitialRow: 0, SongLoopCount: -1, PlayUntilOrder: playUntilOrder, PlayUntilRow: playUntilRow}
	if err := initTick(&m.ticker, &m, ts); err != nil {
		t.Fatalf("initTick error: %v", err)
	}

	err := m.Advance()
	if !errors.Is(err, song.ErrStopSong) {
		t.Fatalf("expected ErrStopSong at play-until position, got %v", err)
	}
}

func TestTickPropagatesAdvanceStopWithoutRender(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		songData: stopSongData{},
	}

	s := sampler.NewSampler(10, 2, 1, nil)

	if err := m.Tick(s); !errors.Is(err, song.ErrStopSong) {
		t.Fatalf("expected ErrStopSong from Advance, got %v", err)
	}
}

func TestTickRendersWhenSequencingContinues(t *testing.T) {
	pat := song.Pattern{song.Row(0), song.Row(1)}
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		globals:  globals[stubGV]{tempo: 2, bpm: 6, gv: stubGV(1), mv: 1},
		ms:       &settings.MachineSettings[stubPeriod, stubGV, stubGV, stubGV, stubPan]{PeriodConverter: stubPeriodCalc{}},
		songData: playSongData{pat: pat},
	}

	called := 0
	s := sampler.NewSampler(10, 2, 1, func(*output.PremixData) {
		called++
	})

	if err := m.Tick(s); err != nil {
		t.Fatalf("unexpected error from Tick: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected OnGenerate to be called once, got %d", called)
	}
}

func TestTickRendersWithNilOnGenerate(t *testing.T) {
	pat := song.Pattern{song.Row(0)}
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		globals:  globals[stubGV]{tempo: 1, bpm: 6, gv: stubGV(1), mv: 1},
		ms:       &settings.MachineSettings[stubPeriod, stubGV, stubGV, stubGV, stubPan]{PeriodConverter: stubPeriodCalc{}},
		songData: playSongData{pat: pat},
	}

	s := sampler.NewSampler(10, 2, 1, nil)

	if err := m.Tick(s); err != nil {
		t.Fatalf("unexpected error from Tick with nil OnGenerate: %v", err)
	}
}

func TestAdvanceSeparatesSequencingFromRendering(t *testing.T) {
	pat := song.Pattern{song.Row(0), song.Row(1)}
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		globals:  globals[stubGV]{tempo: 2},
		songData: seqSongData{pat: pat},
	}

	ch := render.Channel[stubPeriod]{}
	voice := &doneVoice{}
	ch.StartVoice(voice, nil)
	m.actualOutputs = []render.Channel[stubPeriod]{ch}

	if err := initTick(&m.ticker, &m, tickerSettings{InitialOrder: 0, InitialRow: 0, SongLoopCount: -1}); err != nil {
		t.Fatalf("initTick error: %v", err)
	}

	if err := m.Advance(); err != nil {
		t.Fatalf("Advance error: %v", err)
	}

	if voice.ticked != 0 {
		t.Fatalf("expected no rendering during sequencing-only Advance, got %d voice ticks", voice.ticked)
	}
	if m.ticker.current.Tick != 1 || m.ticker.current.Row != 0 || m.ticker.current.Order != 0 {
		t.Fatalf("unexpected sequencing state: %+v", m.ticker.current)
	}
}

func TestRenderDoesNotAdvanceSequencer(t *testing.T) {
	m := machine[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		globals:  globals[stubGV]{tempo: 2, bpm: 6, gv: stubGV(1), mv: 0.5},
		songData: stubSongData{},
	}
	m.ticker.current = Position{Order: 2, Row: 3, Tick: 4}
	m.age = 7

	s := sampler.NewSampler(10, 2, 1, nil)

	prev := m.ticker.current
	if err := m.Render(s); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	if m.age != 7 {
		t.Fatalf("expected age to remain unchanged, got %d", m.age)
	}
	if m.ticker.current != prev {
		t.Fatalf("expected sequencing position unchanged, got %+v", m.ticker.current)
	}
}
