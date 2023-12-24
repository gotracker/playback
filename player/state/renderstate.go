package state

import (
	"time"
)

type RowRenderState struct {
	RenderDetails

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
	return r.RenderDetails.SamplerSpeed
}

func (r RowRenderState) GetDuration() time.Duration {
	return r.RenderDetails.Duration
}

func (r RowRenderState) GetSamples() int {
	return r.RenderDetails.Samples
}
