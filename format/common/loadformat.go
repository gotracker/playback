package common

import (
	"io"

	"github.com/gotracker/playback/player/feature"
)

type ReaderFunc[TSong any] func(r io.Reader, features []feature.Feature) (*TSong, error)

type ManagerFactory[TSong, TManager any] func(*TSong) (*TManager, error)

func Load[TSong, TManager any](r io.Reader, reader ReaderFunc[TSong], factory ManagerFactory[TSong, TManager], features []feature.Feature) (*TManager, error) {
	song, err := reader(r, features)
	if err != nil {
		return nil, err
	}

	m, err := factory(song)

	return m, err
}
