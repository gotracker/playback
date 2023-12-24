package playback

import (
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/period"
)

// ChannelState is the information needed to make an instrument play
type ChannelState[TPeriod period.Period] struct {
	Instrument *instrument.Instrument
	Period     TPeriod
	vol        volume.Volume
	Pos        sampling.Pos
	Pan        panning.Position
}

// Reset sets the render state to defaults
func (s *ChannelState[TPeriod]) Reset() {
	s.Instrument = nil
	var empty TPeriod
	s.Period = empty
	s.Pos = sampling.Pos{}
	s.Pan = panning.CenterAhead
}

func (s *ChannelState[TPeriod]) GetVolume() volume.Volume {
	return s.vol
}

func (s *ChannelState[TPeriod]) SetVolume(vol volume.Volume) {
	if vol != volume.VolumeUseInstVol {
		s.vol = vol
	}
}

func (s *ChannelState[TPeriod]) NoteCut() {
	var empty TPeriod
	s.Period = empty
}
