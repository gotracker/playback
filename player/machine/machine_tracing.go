package machine

import (
	"reflect"

	"github.com/gotracker/playback/index"
	"github.com/heucuva/optional"
)

func trace[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string) {
	m.us.Trace(name)
}

func traceOptionalValueClear[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string, before optional.Value[T]) {
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

func traceWithComment[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name, commentFmt string, commentParams ...any) {
	m.us.TraceWithComment(name, commentFmt, commentParams...)
}

func traceValueChange[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string, before T, after T) {
	if reflect.DeepEqual(before, after) {
		return
	}
	m.us.TraceValueChange(name, before, after)
}

func traceOptionalValueResetWithComment[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string, before optional.Value[T], commentFmt string, commentParams ...any) {
	if v, set := before.Get(); set {
		m.us.TraceValueChangeWithComment(name, v, nil, commentFmt, commentParams...)
	}
}

func traceOptionalValueChangeWithComment[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string, before optional.Value[T], after T, commentFmt string, commentParams ...any) {
	if v, set := before.Get(); set {
		traceValueChangeWithComment(m, name, v, after, commentFmt, commentParams...)
		return
	}

	m.us.TraceValueChangeWithComment(name, nil, after, commentFmt, commentParams...)
}

func traceValueChangeWithComment[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], name string, before, after T, commentFmt string, commentParams ...any) {
	if reflect.DeepEqual(before, after) {
		return
	}

	m.us.TraceValueChangeWithComment(name, before, after, commentFmt, commentParams...)
}

func traceChannel[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name string) {
	m.us.TraceChannel(ch, name)
}

func traceChannelWithComment[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name, commentFmt string, commentParams ...any) {
	m.us.TraceChannelWithComment(ch, name, commentFmt, commentParams...)
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

func traceChannelOptionalValueResetWithComment[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name string, before optional.Value[T], commentFmt string, commentParams ...any) {
	if v, set := before.Get(); set {
		m.us.TraceChannelValueChangeWithComment(ch, name, v, nil, commentFmt, commentParams...)
	}
}

func traceChannelOptionalValueChangeWithComment[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name string, before optional.Value[T], after T, commentFmt string, commentParams ...any) {
	if v, set := before.Get(); set {
		traceChannelValueChangeWithComment(m, ch, name, v, after, commentFmt, commentParams...)
		return
	}

	m.us.TraceChannelValueChangeWithComment(ch, name, nil, after, commentFmt, commentParams...)
}

func traceChannelValueChangeWithComment[T any, TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, name string, before, after T, commentFmt string, commentParams ...any) {
	if reflect.DeepEqual(before, after) {
		return
	}

	m.us.TraceChannelValueChangeWithComment(ch, name, before, after, commentFmt, commentParams...)
}
