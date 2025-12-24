package load

import (
	"testing"

	itPanning "github.com/gotracker/playback/format/it/panning"
)

func TestConvertPanEnvValueScalesToPanRange(t *testing.T) {
	tests := []struct {
		in   int8
		want itPanning.Panning
	}{
		{-32, 0},
		{0, itPanning.DefaultPanning},
		{32, itPanning.MaxPanning},
	}

	for _, tt := range tests {
		got := convertPanEnvValue(tt.in)
		if got != tt.want {
			t.Fatalf("convertPanEnvValue(%d) = %d, want %d", tt.in, got, tt.want)
		}
	}
}
