// Package it does a thing.
package it

import (
	"io"

	"github.com/gotracker/playback/format/common"
	itFeature "github.com/gotracker/playback/format/it/feature"
	"github.com/gotracker/playback/format/it/load"
	itSettings "github.com/gotracker/playback/format/it/settings"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/util"
)

type format struct {
	common.Format
}

var (
	// IT is the exported interface to the IT file loader
	IT = format{}
)

// Load loads an IT file into a playback system
func (f format) Load(filename string, features []feature.Feature) (song.Data, error) {
	r, err := util.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return f.LoadFromReader(r, features)
}

// LoadFromReader loads an IT file on a reader into a playback system
func (format) LoadFromReader(r io.Reader, features []feature.Feature) (song.Data, error) {
	return load.IT(r, features)
}

func (f format) ConvertFeaturesToSettings(us *settings.UserSettings, features []feature.Feature) error {
	for _, feat := range features {
		switch f := feat.(type) {
		case itFeature.LongChannelOutput:
			us.LongChannelOutput = f.Enabled
		case itFeature.NewNoteActions:
			us.EnableNewNoteActions = f.Enabled
		}
	}

	return f.Format.ConvertFeaturesToSettings(us, features)
}

func init() {
	machine.RegisterMachine(itSettings.GetMachineSettings[period.Amiga]())
	machine.RegisterMachine(itSettings.GetMachineSettings[period.Linear]())
}
