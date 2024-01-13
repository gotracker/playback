package common

import (
	"io"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
)

type ReaderFunc func(r io.Reader, features []feature.Feature) (song.Data, error)

type ManagerFactory func(song.Data) (playback.Playback, error)

func Load(r io.Reader, reader ReaderFunc, factory ManagerFactory, features []feature.Feature) (playback.Playback, error) {
	song, err := reader(r, features)
	if err != nil {
		return nil, err
	}

	m, err := factory(song)

	return m, err
}
