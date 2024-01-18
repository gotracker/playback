package format

import (
	"errors"
	"io"
	"os"

	"github.com/gotracker/playback/format/it"
	"github.com/gotracker/playback/format/mod"
	"github.com/gotracker/playback/format/s3m"
	"github.com/gotracker/playback/format/xm"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/song"
)

// Format is an interface to a music file format loader
type Format interface {
	Load(filename string, features []feature.Feature) (song.Data, error)
	LoadFromReader(r io.Reader, features []feature.Feature) (song.Data, error)
	ConvertFeaturesToSettings(us *settings.UserSettings, features []feature.Feature) error
}

var (
	supportedFormats = make(map[string]Format)
)

// Load loads the a file into a playback manager
func Load(filename string, features ...feature.Feature) (song.Data, Format, error) {
	for _, f := range supportedFormats {
		if s, err := f.Load(filename, features); err == nil {
			return s, f, nil
		} else if os.IsNotExist(err) {
			return nil, nil, err
		}
	}
	return nil, nil, errors.New("unsupported format")
}

// LoadFromReader loads a song file on a reader into a playback manager
func LoadFromReader(format string, r io.ReadSeeker, features ...feature.Feature) (song.Data, Format, error) {
	pos, _ := r.Seek(0, io.SeekCurrent)
	if format != "" {
		f, ok := supportedFormats[format]
		if !ok {
			return nil, nil, errors.New("unsupported format")
		}

		_, _ = r.Seek(pos, io.SeekStart)
		if s, err := f.LoadFromReader(r, features); err == nil {
			return s, f, nil
		} else {
			return nil, nil, err
		}
	}

	for _, f := range supportedFormats {
		_, _ = r.Seek(pos, io.SeekStart)
		if s, err := f.LoadFromReader(r, features); err == nil {
			return s, f, nil
		} else if os.IsNotExist(err) {
			return nil, nil, err
		}
	}
	return nil, nil, errors.New("unsupported format")
}

func init() {
	supportedFormats["s3m"] = s3m.S3M
	supportedFormats["mod"] = mod.MOD
	supportedFormats["xm"] = xm.XM
	supportedFormats["it"] = it.IT
}
