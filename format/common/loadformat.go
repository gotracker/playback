package common

import (
	"io"

	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
)

type ReaderFunc func(r io.Reader, features []feature.Feature) (song.Data, error)

func Load(r io.Reader, reader ReaderFunc, features []feature.Feature) (song.Data, error) {
	return reader(r, features)
}
