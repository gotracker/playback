package common

import (
	"io"

	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
)

type ReaderFunc[TSong song.Data] func(r io.Reader, features []feature.Feature) (*TSong, error)

func Load[TSong song.Data](r io.Reader, reader ReaderFunc[TSong], features []feature.Feature) (*TSong, error) {
	return reader(r, features)
}
