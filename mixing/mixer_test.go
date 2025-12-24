package mixing

import (
	"bytes"
	"testing"

	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
)

type fakePanMixer struct {
	matrix volume.Matrix
	calls  int
}

func (p *fakePanMixer) ApplyToMatrix(mtx volume.Matrix) volume.Matrix {
	p.calls++
	return p.matrix.ApplyToMatrix(mtx)
}

func (p *fakePanMixer) Apply(v volume.Volume) volume.Matrix {
	p.calls++
	return p.matrix.Apply(v)
}

func TestMixerFlattenAppliesPanAndVolume(t *testing.T) {
	pan := &fakePanMixer{matrix: volume.Matrix{StaticMatrix: volume.StaticMatrix{1, 0.5}, Channels: 2}}
	flushCalls := 0
	row := []ChannelData{
		{
			{
				Data: MixBuffer{
					{StaticMatrix: volume.StaticMatrix{1, 1}, Channels: 2},
					{StaticMatrix: volume.StaticMatrix{2, 2}, Channels: 2},
				},
				PanMatrix: pan,
				Volume:    0.5,
				Pos:       0,
				Flush:     func() { flushCalls++ },
			},
		},
	}

	m := Mixer{Channels: 2}
	out := m.Flatten(2, row, 1, sampling.Format8BitUnsigned)

	expected := []byte{0xBF, 0x9F, 0xFF, 0xBF}
	if !bytes.Equal(out, expected) {
		t.Fatalf("unexpected mixed data: got %v want %v", out, expected)
	}
	if flushCalls != 1 {
		t.Fatalf("flush should be invoked once, got %d", flushCalls)
	}
	if pan.calls != 1 {
		t.Fatalf("pan Apply should be called once, got %d", pan.calls)
	}
}

func TestMixerFlattenToInts(t *testing.T) {
	pan := &fakePanMixer{matrix: volume.Matrix{StaticMatrix: volume.StaticMatrix{1, 0.5}, Channels: 2}}
	row := []ChannelData{
		{
			{
				Data: MixBuffer{
					{StaticMatrix: volume.StaticMatrix{1, 1}, Channels: 2},
					{StaticMatrix: volume.StaticMatrix{2, 2}, Channels: 2},
				},
				PanMatrix: pan,
				Volume:    0.5,
				Pos:       0,
			},
		},
	}

	m := Mixer{Channels: 2}
	out := m.FlattenToInts(2, 2, 16, row, 1)

	if len(out) != 2 {
		t.Fatalf("expected 2 output channels, got %d", len(out))
	}

	expectedLeft := []int32{16339, 32678}
	expectedRight := []int32{8169, 16339}

	if !bytes.Equal(int32ToBytes(out[0]), int32ToBytes(expectedLeft)) {
		t.Fatalf("unexpected left channel: got %v want %v", out[0], expectedLeft)
	}
	if !bytes.Equal(int32ToBytes(out[1]), int32ToBytes(expectedRight)) {
		t.Fatalf("unexpected right channel: got %v want %v", out[1], expectedRight)
	}
}

func TestMixerFlattenToBuffers(t *testing.T) {
	pan := &fakePanMixer{matrix: volume.Matrix{StaticMatrix: volume.StaticMatrix{1, 0.5}, Channels: 2}}
	row := []ChannelData{
		{
			{
				Data: MixBuffer{
					{StaticMatrix: volume.StaticMatrix{1, 1}, Channels: 2},
					{StaticMatrix: volume.StaticMatrix{2, 2}, Channels: 2},
				},
				PanMatrix: pan,
				Volume:    0.5,
				Pos:       0,
			},
		},
	}

	bufA := make([]byte, 2)
	bufB := make([]byte, 2)

	m := Mixer{Channels: 2}
	m.FlattenTo([][]byte{bufA, bufB}, 2, 2, row, 1, sampling.Format8BitUnsigned)

	expectedA := []byte{0xBF, 0x9F}
	expectedB := []byte{0xFF, 0xBF}

	if !bytes.Equal(bufA, expectedA) {
		t.Fatalf("unexpected first buffer: got %v want %v", bufA, expectedA)
	}
	if !bytes.Equal(bufB, expectedB) {
		t.Fatalf("unexpected second buffer: got %v want %v", bufB, expectedB)
	}
}

// int32ToBytes provides a stable comparison helper without float rounding differences.
func int32ToBytes(in []int32) []byte {
	buf := bytes.Buffer{}
	for _, v := range in {
		// big endian is fine; we just need byte stability for comparison
		buf.Write([]byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)})
	}
	return buf.Bytes()
}
