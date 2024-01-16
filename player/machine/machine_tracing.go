package machine

import (
	"reflect"

	"github.com/gotracker/playback/index"
	"github.com/heucuva/optional"
)

func trace[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string) {
	m.us.Trace(name)
}

func traceOptionalValueClear[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string, before optional.Value[T], comment string) {
	if v, set := before.Get(); set {
		m.us.TraceValueChange(name, v, nil)
	}
}

func traceOptionalValueChange[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string, before optional.Value[T], after T) {
	if v, set := before.Get(); set {
		traceValueChange(m, name, v, after)
		return
	}

	m.us.TraceValueChange(name, nil, after)
}

func traceWithComment[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name, comment string) {
	m.us.TraceWithComment(name, comment)
}

func traceValueChange[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string, before T, after T) {
	if reflect.DeepEqual(before, after) {
		return
	}
	m.us.TraceValueChange(name, before, after)
}

func traceOptionalValueResetWithComment[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string, before optional.Value[T], comment string) {
	if v, set := before.Get(); set {
		m.us.TraceValueChangeWithComment(name, v, nil, comment)
	}
}

func traceOptionalValueChangeWithComment[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string, before optional.Value[T], after T, comment string) {
	if v, set := before.Get(); set {
		traceValueChangeWithComment(m, name, v, after, comment)
		return
	}

	m.us.TraceValueChangeWithComment(name, nil, after, comment)
}

func traceValueChangeWithComment[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string, before, after T, comment string) {
	if reflect.DeepEqual(before, after) {
		return
	}

	m.us.TraceValueChangeWithComment(name, before, after, comment)
}

func traceChannel[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name string) {
	m.us.TraceChannel(ch, name)
}

func traceChannelWithComment[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name, comment string) {
	m.us.TraceChannelWithComment(ch, name, comment)
}

func traceChannelOptionalValueReset[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name string, before optional.Value[T]) {
	if v, set := before.Get(); set {
		m.us.TraceChannelValueChange(ch, name, v, nil)
	}
}

func traceChannelOptionalValueChange[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name string, before optional.Value[T], after T) {
	if v, set := before.Get(); set {
		traceChannelValueChange(m, ch, name, v, after)
		return
	}

	m.us.TraceChannelValueChange(ch, name, nil, after)
}

func traceChannelValueChange[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name string, before, after T) {
	if reflect.DeepEqual(before, after) {
		return
	}

	m.us.TraceChannelValueChange(ch, name, before, after)
}

func traceChannelOptionalValueResetWithComment[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name string, before optional.Value[T], comment string) {
	if v, set := before.Get(); set {
		m.us.TraceChannelValueChangeWithComment(ch, name, v, nil, comment)
	}
}

func traceChannelOptionalValueChangeWithComment[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name string, before optional.Value[T], after T, comment string) {
	if v, set := before.Get(); set {
		traceChannelValueChangeWithComment(m, ch, name, v, after, comment)
		return
	}

	m.us.TraceChannelValueChangeWithComment(ch, name, nil, after, comment)
}

func traceChannelValueChangeWithComment[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name string, before, after T, comment string) {
	if reflect.DeepEqual(before, after) {
		return
	}

	m.us.TraceChannelValueChangeWithComment(ch, name, before, after, comment)
}
