package mod

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m"
	"github.com/gotracker/playback/format/settings"
)

type format struct{}

var (
	// MOD is the exported interface to the MOD file loader
	MOD = format{}
)

// Load loads an MOD file into the song state
func (f format) Load(filename string, s *settings.Settings) (playback.Playback, error) {
	// we really just load the mod into an S3M layout, since S3M is essentially a superset
	return s3m.LoadMOD(filename, s)
}
