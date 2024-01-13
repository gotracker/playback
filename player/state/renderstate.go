package state

import (
	"time"

	"github.com/gotracker/playback/player/state/render"
)

type RowRenderState struct {
	render.Details

	TicksThisRow int
	CurrentTick  int
}

func (r RowRenderState) GetTicksThisRow() int {
	return r.TicksThisRow
}

func (r RowRenderState) GetCurrentTick() int {
	return r.CurrentTick
}

func (r RowRenderState) GetSamplerSpeed() float32 {
	return r.Details.SamplerSpeed
}

func (r RowRenderState) GetDuration() time.Duration {
	return r.Details.Duration
}

func (r RowRenderState) GetSamples() int {
	return r.Details.Samples
}
