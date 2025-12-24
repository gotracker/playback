package playback

import (
	"github.com/gotracker/playback/mixing/sampling"

	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/voice/types"
)

// ChannelState is the information needed to make an instrument play
type ChannelState[TPeriod types.Period, TVolume types.Volume, TPanning types.Panning] struct {
	inst   instrument.InstrumentIntf
	period TPeriod
	vol    TVolume
	pos    sampling.Pos
	pan    TPanning
}

// Reset sets the render state to defaults
func (s *ChannelState[TPeriod, TVolume, TPanning]) Reset() {
	s.inst = nil
	var emptyPeriod TPeriod
	s.period = emptyPeriod
	s.pos = sampling.Pos{}
	var emptyPan TPanning
	s.pan = emptyPan
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
	s.period = empty
}

func (s *ChannelState[TPeriod, TVolume, TPanning]) Instrument() instrument.InstrumentIntf {
	return s.inst
}

func (s *ChannelState[TPeriod, TVolume, TPanning]) SetInstrument(inst instrument.InstrumentIntf) {
	s.inst = inst
}

func (s *ChannelState[TPeriod, TVolume, TPanning]) Period() TPeriod {
	return s.period
}

func (s *ChannelState[TPeriod, TVolume, TPanning]) SetPeriod(p TPeriod) {
	s.period = p
}

func (s *ChannelState[TPeriod, TVolume, TPanning]) Pos() sampling.Pos {
	return s.pos
}

func (s *ChannelState[TPeriod, TVolume, TPanning]) SetPos(pos sampling.Pos) {
	s.pos = pos
}

func (s *ChannelState[TPeriod, TVolume, TPanning]) Pan() TPanning {
	return s.pan
}

func (s *ChannelState[TPeriod, TVolume, TPanning]) SetPan(p TPanning) {
	s.pan = p
}
