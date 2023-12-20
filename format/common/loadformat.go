package common

import (
	"io"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/player/feature"
)

type ReaderFunc[TSong any] func(r io.Reader, features []feature.Feature) (*TSong, error)

type ManagerFactory[TSong any] func(*TSong) (playback.Playback, error)

func Load[TSong any](r io.Reader, reader ReaderFunc[TSong], factory ManagerFactory[TSong], features []feature.Feature) (playback.Playback, error) {
	song, err := reader(r, features)
	if err != nil {
		return nil, err
	}

	m, err := factory(song)

	return m, err
}
