package channel

import (
	"github.com/gotracker/playback/memory"
	"github.com/gotracker/playback/tremor"
)

// Memory is the storage object for custom effect/command values
type Memory struct {
	porta         memory.Value[DataEffect]
	vibratoSpeed  memory.Value[DataEffect]
	vibratoDepth  memory.Value[DataEffect]
	tremoloSpeed  memory.Value[DataEffect]
	tremoloDepth  memory.Value[DataEffect]
	sampleOffset  memory.Value[DataEffect]
	tempoDecrease memory.Value[DataEffect]
	tempoIncrease memory.Value[DataEffect]
	lastNonZero   memory.Value[DataEffect]

	tremorMem tremor.Tremor

	Shared *SharedMemory
}

// Porta gets or sets the most recent non-zero value (or input) for any Portamento command
func (m *Memory) Porta(input DataEffect) DataEffect {
	return m.porta.Coalesce(input)
}

// Vibrato gets or sets the most recent non-zero value (or input) for Vibrato
func (m *Memory) Vibrato(input DataEffect) (DataEffect, DataEffect) {
	// vibrato is unusual, because each nibble is treated uniquely
	vx := m.vibratoSpeed.Coalesce(input >> 4)
	vy := m.vibratoDepth.Coalesce(input & 0x0f)
	return vx, vy
}

// Tremolo gets or sets the most recent non-zero value (or input) for Vibrato
func (m *Memory) Tremolo(input DataEffect) (DataEffect, DataEffect) {
	// tremolo is unusual, because each nibble is treated uniquely
	vx := m.tremoloSpeed.Coalesce(input >> 4)
	vy := m.tremoloDepth.Coalesce(input & 0x0f)
	return vx, vy
}

// SampleOffset gets or sets the most recent non-zero value (or input) for Sample Offset
func (m *Memory) SampleOffset(input DataEffect) DataEffect {
	return m.sampleOffset.Coalesce(input)
}

// TempoDecrease gets or sets the most recent non-zero value (or input) for Tempo Decrease
func (m *Memory) TempoDecrease(input DataEffect) DataEffect {
	return m.tempoDecrease.Coalesce(input)
}

// TempoIncrease gets or sets the most recent non-zero value (or input) for Tempo Increase
func (m *Memory) TempoIncrease(input DataEffect) DataEffect {
	return m.tempoIncrease.Coalesce(input)
}

// LastNonZero gets or sets the most recent non-zero value (or input)
func (m *Memory) LastNonZero(input DataEffect) DataEffect {
	return m.lastNonZero.Coalesce(input)
}

// LastNonZero gets or sets the most recent non-zero value (or input)
func (m *Memory) LastNonZeroXY(input DataEffect) (DataEffect, DataEffect) {
	xy := m.LastNonZero(input)
	return xy >> 4, xy & 0x0f
}

// TremorMem returns the Tremor object
func (m *Memory) TremorMem() *tremor.Tremor {
	return &m.tremorMem
}

// Retrigger is called when a voice is triggered
func (m *Memory) Retrigger() {
}

// StartOrder is called when the first order's row at tick 0 is started
func (m *Memory) StartOrder0() {
	if m.Shared.ResetMemoryAtStartOfOrder0 {
		m.porta.Reset()
		m.vibratoSpeed.Reset()
		m.vibratoDepth.Reset()
		m.tremoloSpeed.Reset()
		m.tremoloDepth.Reset()
		m.sampleOffset.Reset()
		m.tempoDecrease.Reset()
		m.tempoIncrease.Reset()
		m.lastNonZero.Reset()
	}
}
