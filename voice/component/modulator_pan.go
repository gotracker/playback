package component

import (
	"fmt"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/types"
)

// PanModulator is an pan (spatial) modulator
type PanModulator[TPanning types.Panning] struct {
	settings PanModulatorSettings[TPanning]
	pan      TPanning
	delta    types.PanDelta
	final    TPanning
}

type PanModulatorSettings[TPanning types.Panning] struct {
	InitialPan TPanning
}

func (p *PanModulator[TPanning]) Setup(settings PanModulatorSettings[TPanning]) {
	p.settings = settings
	p.pan = settings.InitialPan
	p.delta = 0
	p.updateFinal()
}

func (p PanModulator[TPanning]) Clone() PanModulator[TPanning] {
	m := p
	return m
}

// SetPan sets the current panning
func (p *PanModulator[TPanning]) SetPan(pan TPanning) {
	p.pan = pan
	p.updateFinal()
}

// GetPan returns the current panning
func (p PanModulator[TPanning]) GetPan() TPanning {
	return p.pan
}

// SetPanDelta sets the current panning delta
func (p *PanModulator[TPanning]) SetPanDelta(d types.PanDelta) {
	p.delta = d
	p.updateFinal()
}

// GetPanDelta returns the current panning delta
func (p PanModulator[TPanning]) GetPanDelta() types.PanDelta {
	return p.delta
}

// GetFinalPan returns the current panning
func (p PanModulator[TPanning]) GetFinalPan() TPanning {
	return p.final
}

func (p PanModulator[TPanning]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("pan{%v} delta{%v}",
		p.pan,
		p.delta,
	), comment)
}

func (p *PanModulator[TPanning]) updateFinal() {
	p.final = types.AddPanningDelta(p.pan, p.delta)
}
