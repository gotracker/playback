package playback

import (
	"github.com/gotracker/gomixing/sampling"

	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/voice/types"
)

// ChannelState is the information needed to make an instrument play
type ChannelState[TPeriod types.Period, TVolume types.Volume, TPanning types.Panning] struct {
	Instrument instrument.InstrumentIntf
	Period     TPeriod
	vol        TVolume
	Pos        sampling.Pos
	Pan        TPanning
}

// Reset sets the render state to defaults
func (s *ChannelState[TPeriod, TVolume, TPanning]) Reset() {
	s.Instrument = nil
	var emptyPeriod TPeriod
	s.Period = emptyPeriod
	s.Pos = sampling.Pos{}
	var emptyPan TPanning
	s.Pan = emptyPan
}

func (s *ChannelState[TPeriod, TVolume, TPanning]) GetVolume() TVolume {
	return s.vol
}

func (s *ChannelState[TPeriod, TVolume, TPanning]) SetVolume(vol TVolume) {
	if !vol.IsUseInstrumentVol() {
		s.vol = vol
	}
}

func (s *ChannelState[TPeriod, TVolume, TPanning]) NoteCut() {
	var empty TPeriod
	s.Period = empty
}
