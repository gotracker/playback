package tracing

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player"
	"github.com/gotracker/playback/song"
)

type Tracing[TPeriod period.Period, TMemory any, TData song.ChannelData] struct {
	playback.Playback
	ChannelGetter func(c int) playback.Channel[TPeriod, TMemory, TData]
}

func (m Tracing[TPeriod, TMemory, TData]) OutputTraces(out chan<- func(w io.Writer)) {
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

func (m Tracing[TPeriod, TMemory, TData]) outputGlobalTrace() func(w io.Writer) {
	gs := player.NewTracingTable("=== global ===",
		"globalVolume",
		"mixerVolume",
		"currentOrder",
		"currentRow",
	)
	gs.AddRow(
		m.Playback.GetGlobalVolume(),
		m.Playback.GetMixerVolume(),
		m.Playback.GetCurrentOrder(),
		m.Playback.GetCurrentRow(),
	)

	return func(w io.Writer) {
		fmt.Fprintln(w)

		tw := tabwriter.NewWriter(w, 1, 1, 1, ' ', 0)
		defer tw.Flush()

		gs.Fprintln(tw, "\t", false)
	}
}

func (m Tracing[TPeriod, TMemory, TData]) outputRenderTrace() func(w io.Writer) {
	r := m.Playback.GetRenderState()
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
		fmt.Sprint(r.GetSamplerSpeed()),
		fmt.Sprint(r.GetDuration()),
		fmt.Sprint(r.GetSamples()),
		fmt.Sprint(r.GetTicksThisRow()),
		fmt.Sprint(r.GetCurrentTick()),
	)

	return func(w io.Writer) {
		fmt.Fprintln(w)

		tw := tabwriter.NewWriter(w, 1, 1, 1, ' ', 0)
		defer tw.Flush()

		rs.Fprintln(tw, "\t", false)
	}
}

func (m Tracing[TPeriod, TMemory, TData]) outputChannelsTrace() func(w io.Writer) {
	cs := player.NewTracingTable("=== channels ===",
		append(
			[]string{
				"Channel",
				"ChannelVolume",
				"ActiveEffect",
				"TrackData",
				"RetriggerCount",
				"Semitone",
				"UseTargetPeriod",
				"NewNoteAction",
			},
			ChannelStateHeaders()...,
		)...,
	)

	for c := 0; c < m.Playback.GetNumChannels(); c++ {
		ch := m.ChannelGetter(c)
		if ch == nil {
			continue
		}
		var trackData string
		effects := ch.GetActiveEffects()
		if len(effects) == 0 {
			effects = []playback.Effect{nil}
		}
		trackData = fmt.Sprint(ch.GetChannelData())
		var activeEffect []string
		for _, effect := range effects {
			if effect != nil {
				effectTypes := playback.GetEffectNames(effect)
				activeEffect = append(activeEffect, strings.Join(effectTypes, ","))
			}
		}

		cs.AddRow(
			append(
				[]any{
					c + 1,
					ch.GetChannelVolume(),
					strings.Join(activeEffect, ","),
					trackData,
					ch.GetRetriggerCount(),
					ch.GetNoteSemitone(),
					ch.GetUseTargetPeriod(),
					ch.GetNewNoteAction(),
				},
				ChannelState[TPeriod](ch.GetActiveState())...,
			)...,
		)
	}

	return func(w io.Writer) {
		fmt.Fprintln(w)

		tw := tabwriter.NewWriter(w, 1, 1, 1, ' ', 0)
		defer tw.Flush()

		cs.Fprintln(tw, "\t", true)
	}
}
