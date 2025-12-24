package filter

import (
	"math"
	"testing"

	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/mixing/volume"
)

func almostEqualVol(a, b volume.Volume, tol float64) bool {
	return math.Abs(float64(a-b)) <= tol
}

// helper to compare two matrices of equal channel counts within tolerance
func assertMatrixAlmostEqual(t *testing.T, got, want volume.Matrix, tol float64) {
	t.Helper()
	if got.Channels != want.Channels {
		t.Fatalf("channel mismatch: got %d want %d", got.Channels, want.Channels)
	}
	for i := 0; i < got.Channels; i++ {
		if !almostEqualVol(got.StaticMatrix[i], want.StaticMatrix[i], tol) {
			t.Fatalf("channel %d mismatch: got %v want %v", i, got.StaticMatrix[i], want.StaticMatrix[i])
		}
	}
}

func TestAmigaLPFFilterProgression(t *testing.T) {
	f := NewAmigaLPF(frequency.Frequency(6550))

	dry := volume.Matrix{StaticMatrix: volume.StaticMatrix{1}, Channels: 1}

	first := f.Filter(dry)
	expectedFirst := dry.StaticMatrix[0] * f.a0
	if !almostEqualVol(first.StaticMatrix[0], expectedFirst, 1e-6) {
		t.Fatalf("first sample mismatch: got %v want %v", first.StaticMatrix[0], expectedFirst)
	}

	second := f.Filter(dry)
	expectedSecond := dry.StaticMatrix[0]*f.a0 + expectedFirst*f.b0
	if !almostEqualVol(second.StaticMatrix[0], expectedSecond, 1e-6) {
		t.Fatalf("second sample mismatch: got %v want %v", second.StaticMatrix[0], expectedSecond)
	}
}

func TestEchoFilterZeroDelayNoFeedback(t *testing.T) {
	echo := EchoFilter{
		EchoFilterSettings: EchoFilterSettings{
			WetDryMix:  0.5,
			Feedback:   0,
			LeftDelay:  0.125, // yields 1-sample delay at 4Hz playback rate
			RightDelay: 0.125,
			PanDelay:   0,
		},
	}
	echo.SetPlaybackRate(frequency.Frequency(4))

	dry := volume.Matrix{StaticMatrix: volume.StaticMatrix{1, -1}, Channels: 2}

	first := echo.Filter(dry)
	expectedFirst := dry.Apply(volume.Volume(0.5))
	assertMatrixAlmostEqual(t, first, expectedFirst, 1e-6)

	second := echo.Filter(dry)
	assertMatrixAlmostEqual(t, second, dry, 1e-6)
}

func TestResonantFilterBypassWhenWideOpen(t *testing.T) {
	rf := NewITResonantFilter(0xFF, 0x00, false, false)
	rf.SetPlaybackRate(frequency.Frequency(44100))

	dry := volume.Matrix{StaticMatrix: volume.StaticMatrix{0.25, -0.25}, Channels: 2}

	wet := rf.Filter(dry)
	if wet != dry {
		t.Fatalf("expected bypass output to equal input: got %v want %v", wet, dry)
	}
	if rf.(*ResonantFilter).enabled {
		t.Fatalf("filter should be disabled for wide-open cutoff and zero resonance")
	}
}

func TestResonantFilterAppliesCoefficients(t *testing.T) {
	rf := NewITResonantFilter(0x80|64, 0x80|32, false, false)
	rf.SetPlaybackRate(frequency.Frequency(48000))

	dry := volume.Matrix{StaticMatrix: volume.StaticMatrix{1}, Channels: 1}

	wet := rf.Filter(dry)
	expected := dry.StaticMatrix[0] * rf.(*ResonantFilter).a0
	if !almostEqualVol(wet.StaticMatrix[0], expected, 1e-5) {
		t.Fatalf("expected filtered output near %v, got %v", expected, wet.StaticMatrix[0])
	}
}
