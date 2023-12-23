package playback

import (
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/player"
)

func (m *manager[TPeriod]) OutputTraces(out chan<- func(w io.Writer)) {
	outputs := []func(w io.Writer){
		m.outputGlobalTrace(),
		m.outputRenderTrace(),
		m.outputChannelsTrace(),
	}
	out <- func(w io.Writer) {
		fmt.Fprintln(w, "################################################")
		for _, fn := range outputs {
			fn(w)
		}

		fmt.Fprintln(w)
	}
}

func (m *manager[TPeriod]) outputGlobalTrace() func(w io.Writer) {
	gs := player.NewTracingTable("=== global ===",
		"globalVolume",
		"mixerVolume",
		"currentOrder",
		"currentRow",
	)
	gs.AddRow(
		m.GetGlobalVolume(),
		m.GetMixerVolume(),
		m.GetCurrentOrder(),
		m.GetCurrentRow(),
	)

	return func(w io.Writer) {
		fmt.Fprintln(w)

		tw := tabwriter.NewWriter(w, 1, 1, 1, ' ', 0)
		defer tw.Flush()

		gs.Fprintln(tw, "\t", false)
	}
}

func (m *manager[TPeriod]) outputRenderTrace() func(w io.Writer) {
	r := m.rowRenderState
	if r == nil {
		return func(w io.Writer) {}
	}

	rs := player.NewTracingTable("=== rowRenderState ===",
		"samplerSpeed",
		"tickDuration",
		"samplesPerTick",
		"ticksThisRow",
		"currentTick",
	)
	rs.AddRow(
		fmt.Sprint(r.SamplerSpeed),
		fmt.Sprint(r.Duration),
		fmt.Sprint(r.Samples),
		fmt.Sprint(r.ticksThisRow),
		fmt.Sprint(r.currentTick),
	)

	return func(w io.Writer) {
		fmt.Fprintln(w)

		tw := tabwriter.NewWriter(w, 1, 1, 1, ' ', 0)
		defer tw.Flush()

		rs.Fprintln(tw, "\t", false)
	}
}

func (m *manager[TPeriod]) outputChannelsTrace() func(w io.Writer) {
	cs := player.NewTracingTable("=== channels ===",
		"Channel",
		"ChannelVolume",
		"ActiveEffect",
		"ActiveEffectType",
		"TrackData",
		"RetriggerCount",
		"Semitone",
		"UseTargetPeriod",
		"PanEnabled",
		"NewNoteAction",
	)

	for c, ch := range m.channels {
		var trackData string
		effects := ch.GetActiveEffects()
		if len(effects) == 0 {
			effects = []playback.Effect{nil}
		}
		trackData = fmt.Sprint(ch.GetChannelData())
		for _, effect := range effects {
			var (
				activeEffect     string
				activeEffectType string
			)
			if effect != nil {
				activeEffect = fmt.Sprint(effect)
				activeEffectType = reflect.TypeOf(effect).Name()
			}
			cs.AddRow(
				c+1,
				ch.GetChannelVolume(),
				activeEffect,
				activeEffectType,
				trackData,
				ch.RetriggerCount,
				ch.Semitone,
				ch.UseTargetPeriod,
				ch.PanEnabled,
				ch.NewNoteAction,
			)
		}
	}

	return func(w io.Writer) {
		fmt.Fprintln(w)

		tw := tabwriter.NewWriter(w, 1, 1, 1, ' ', 0)
		defer tw.Flush()

		cs.Fprintln(tw, "\t", true)
	}
}
