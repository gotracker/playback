package oscillator

// WaveTableSelect is the selection code for which waveform to use in an oscillator
type WaveTableSelect uint8

// Oscillator is an oscillator
type Oscillator interface {
	Clone() Oscillator
	GetWave(depth float32) float32
	Advance(speed int)
	SetWaveform(table WaveTableSelect)
	Reset(hard ...bool)
}
