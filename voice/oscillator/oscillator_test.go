package oscillator

import "testing"

type fakeOsc struct {
	wave  WaveTableSelect
	adv   int
	reset int
}

func (f *fakeOsc) Clone() Oscillator                 { cp := *f; return &cp }
func (f *fakeOsc) GetWave(depth float32) float32     { return float32(f.wave) + depth }
func (f *fakeOsc) Advance(speed int)                 { f.adv += speed }
func (f *fakeOsc) SetWaveform(table WaveTableSelect) { f.wave = table }
func (f *fakeOsc) GetWaveform() WaveTableSelect      { return f.wave }
func (f *fakeOsc) HardReset()                        { f.reset++ }
func (f *fakeOsc) Reset()                            { f.reset++ }

func TestOscillatorInterfaceMethods(t *testing.T) {
	f := &fakeOsc{wave: 2}
	clone := f.Clone().(*fakeOsc)
	clone.SetWaveform(3)
	clone.Advance(4)
	clone.Reset()

	if f.wave != 2 {
		t.Fatalf("expected original waveform unchanged, got %d", f.wave)
	}
	if clone.wave != 3 || clone.adv != 4 || clone.reset != 1 {
		t.Fatalf("unexpected clone state: %+v", clone)
	}
	if clone.GetWave(1) != 4 {
		t.Fatalf("GetWave unexpected result")
	}
}
