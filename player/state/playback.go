package state

import (
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/period"
)

// Playback is the information needed to make an instrument play
type Playback[TPeriod period.Period] struct {
	Instrument *instrument.Instrument
	Period     *TPeriod
	Volume     volume.Volume
	Pos        sampling.Pos
	Pan        panning.Position
}

// Reset sets the render state to defaults
func (p *Playback[TPeriod]) Reset() {
	p.Instrument = nil
	p.Period = nil
	p.Pos = sampling.Pos{}
	p.Pan = panning.CenterAhead
}
