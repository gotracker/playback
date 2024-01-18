package mod

import (
	"io"

	"github.com/gotracker/playback/format/common"
	"github.com/gotracker/playback/format/s3m/load"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/util"
)

type format struct {
	common.Format
}

var (
	// MOD is the exported interface to the MOD file loader
	MOD = format{}
)

// Load loads an MOD file into a playback system
func (f format) Load(filename string, features []feature.Feature) (song.Data, error) {
	r, err := util.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return f.LoadFromReader(r, features)
}

// LoadFromReader loads a MOD file on a reader into a playback system
func (format) LoadFromReader(r io.Reader, features []feature.Feature) (song.Data, error) {
	// we really just load the mod into an S3M layout, since S3M is essentially a superset
	return load.MOD(r, features)
}
