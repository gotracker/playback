package song

import (
	"testing"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/mixing/volume"
)

type testVol int

func (testVol) IsInvalid() bool           { return false }
func (testVol) IsUseInstrumentVol() bool  { return false }
func (v testVol) ToVolume() volume.Volume { return volume.Volume(v) }

type stubRow[T Volume] struct{ n int }

func (s stubRow[T]) Len() int { return s.n }
func (s stubRow[T]) ForEach(fn func(ch index.Channel, d ChannelData[T]) (bool, error)) error {
	for i := 0; i < s.n; i++ {
		if stop, err := fn(index.Channel(i), nil); stop || err != nil {
			return err
		}
	}
	return nil
}

func TestGetRowNumChannels(t *testing.T) {
	if GetRowNumChannels(stubRow[testVol]{n: 3}) != 3 {
		t.Fatalf("expected 3 channels")
	}
	if GetRowNumChannels(struct{}{}) != 0 {
		t.Fatalf("expected 0 for unknown type")
	}
}

func TestForEachRowChannel(t *testing.T) {
	r := stubRow[testVol]{n: 2}
	count := 0
	err := ForEachRowChannel[testVol](r, func(ch index.Channel, d ChannelData[testVol]) (bool, error) {
		count++
		return false, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected to iterate 2 channels, got %d", count)
	}
}
