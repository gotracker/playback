package voice

import (
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
)

func (v xmVoice[TPeriod]) DumpState(ch index.Channel, t tracing.Tracer) {
	if t == nil {
		return
	}

	v.KeyModulator.DumpState(ch, t, "xmVoice.KeyModulator")
	if v.voicer != nil {
		v.voicer.DumpState(ch, t, "xmVoice.voicer")
	} else {
		t.TraceChannelWithComment(ch, "nil", "xmVoice.voicer")
	}
	v.amp.DumpState(ch, t, "xmVoice.amp")
	v.freq.DumpState(ch, t, "xmVoice.freq")
	v.pan.DumpState(ch, t, "xmVoice.pan")
	v.volEnv.DumpState(ch, t, "xmVoice.volEnv")
	v.panEnv.DumpState(ch, t, "xmVoice.panEnv")
	v.vol0Opt.DumpState(ch, t, "xmVoice.vol0Opt")
	//voiceFilter
	//pluginFilter
}
