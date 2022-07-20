package common

import (
	"io"

	"github.com/gotracker/playback/settings"
)

type ReaderFunc[TSong any] func(r io.Reader, s *settings.Settings) (*TSong, error)

type ManagerFactory[TSong, TManager any] func(*TSong) (*TManager, error)

func Load[TSong, TManager any](r io.Reader, reader ReaderFunc[TSong], factory ManagerFactory[TSong, TManager], s *settings.Settings) (*TManager, error) {
	song, err := reader(r, s)
	if err != nil {
		return nil, err
	}

	m, err := factory(song)

	return m, err
}
