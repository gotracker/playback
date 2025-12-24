package common

import (
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/machine/settings"
)

type Format struct{}

func (Format) ConvertFeaturesToSettings(us *settings.UserSettings, features []feature.Feature) error {
	for _, feat := range features {
		switch f := feat.(type) {
		case feature.SongLoop:
			us.SongLoopCount = f.Count
		case feature.StartOrderAndRow:
			if o, set := f.Order.Get(); set {
				us.Start.Order.Set(index.Order(o))
			}
			if r, set := f.Row.Get(); set {
				us.Start.Row.Set(index.Row(r))
			}
		case feature.PlayUntilOrderAndRow:
			us.PlayUntil.Order.Set(index.Order(f.Order))
			us.PlayUntil.Row.Set(index.Row(f.Row))
		case feature.SetDefaultTempo:
			us.Start.Tempo = f.Tempo
		case feature.SetDefaultBPM:
			us.Start.BPM = f.BPM
		case feature.IgnoreUnknownEffect:
			us.IgnoreUnknownEffect = f.Enabled
		case feature.QuirksMode:
			if prof, ok := f.Profile.Get(); ok {
				us.Quirks.Profile.Set(prof)
			}
			if linear, ok := f.LinearSlides.Get(); ok {
				us.Quirks.LinearSlidesOverride.Set(linear)
			}
		}
	}
	return nil
}
