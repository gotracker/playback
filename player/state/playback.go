package state

import (
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/period"
)

// Playback is the information needed to make an instrument play
type Playback struct {
	Instrument *instrument.Instrument
	Period     period.Period
	Volume     volume.Volume
	Pos        sampling.Pos
	Pan        panning.Position
}

// Reset sets the render state to defaults
func (p *Playback) Reset() {
	p.Instrument = nil
	p.Period = nil
	p.Pos = sampling.Pos{}
	p.Pan = panning.CenterAhead
}
