package voice

import (
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
)

func (v itVoice[TPeriod]) DumpState(ch index.Channel, t tracing.Tracer) {
	if t == nil {
		return
	}

	v.KeyModulator.DumpState(ch, t, "itVoice.KeyModulator")
	if v.voicer != nil {
		v.voicer.DumpState(ch, t, "itVoice.voicer")
	} else {
		t.TraceChannelWithComment(ch, "nil", "itVoice.voicer")
	}
	v.amp.DumpState(ch, t, "itVoice.amp")
	v.freq.DumpState(ch, t, "itVoice.freq")
	v.pan.DumpState(ch, t, "itVoice.pan")
	v.volEnv.DumpState(ch, t, "itVoice.volEnv")
	v.pitchEnv.DumpState(ch, t, "itVoice.pitchEnv")
	v.panEnv.DumpState(ch, t, "itVoice.panEnv")
	v.filterEnv.DumpState(ch, t, "itVoice.filterEnv")
	v.vol0Opt.DumpState(ch, t, "itVoice.vol0Opt")
	//voiceFilter
	//pluginFilter
}
