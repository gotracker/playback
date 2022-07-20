package format

import (
	"errors"
	"os"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it"
	"github.com/gotracker/playback/format/mod"
	"github.com/gotracker/playback/format/s3m"
	"github.com/gotracker/playback/format/settings"
	"github.com/gotracker/playback/format/xm"
	"github.com/gotracker/playback/song"
)

var (
	supportedFormats = make(map[string]playback.Format[song.ChannelData])
)

// Load loads the a file into a playback manager
func Load(filename string, options ...settings.OptionFunc) (playback.Playback, playback.Format[song.ChannelData], error) {
	s := &settings.Settings{}
	for _, opt := range options {
		if err := opt(s); err != nil {
			return nil, nil, err
		}
	}

	for _, f := range supportedFormats {
		if playback, err := f.Load(filename, s); err == nil {
			return playback, f, nil
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
