package render

import (
	"testing"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/period"
)

type stubFilter struct {
	factor       float32
	playbackRate float32
}

func (s *stubFilter) Filter(m volume.Matrix) volume.Matrix {
	for i := 0; i < m.Channels; i++ {
		m.StaticMatrix[i] = volume.Volume(float32(m.StaticMatrix[i]) * s.factor)
	}
	return m
}
func (s *stubFilter) SetPlaybackRate(pr frequency.Frequency) { s.playbackRate = float32(pr) }
func (s *stubFilter) UpdateEnv(uint8)                        {}
func (s *stubFilter) Clone() filter.Filter                   { c := *s; return &c }

func TestChannelApplyFilterPipeline(t *testing.T) {
	plugin := &stubFilter{factor: 2}
	output := &stubFilter{factor: 0.5}
	ch := Channel[period.Linear]{
		PluginFilter: plugin,
		OutputFilter: output,
		GlobalVolume: 0.5,
	}

	dry := volume.Matrix{StaticMatrix: volume.StaticMatrix{1}, Channels: 1}
	wet := ch.ApplyFilter(dry)
	// dry -> plugin *2 -> 2; apply GlobalVolume 0.5 -> 1; output *0.5 -> 0.5
	if wet.StaticMatrix[0] != 0.5 {
		t.Fatalf("unexpected filtered value: %v", wet.StaticMatrix[0])
	}
}
