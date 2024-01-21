package voice

import (
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
)

func (v s3mVoice) DumpState(ch index.Channel, t tracing.Tracer) {
	if t == nil {
		return
	}

	v.KeyModulator.DumpState(ch, t, "s3mVoice.KeyModulator")
	if v.voicer != nil {
		v.voicer.DumpState(ch, t, "s3mVoice.voicer")
	} else {
		t.TraceChannelWithComment(ch, "nil", "s3mVoice.voicer")
	}
	v.AmpModulator.DumpState(ch, t, "s3mVoice.amp")
	v.FreqModulator.DumpState(ch, t, "s3mVoice.freq")
	v.PanModulator.DumpState(ch, t, "s3mVoice.pan")
	v.vol0Opt.DumpState(ch, t, "s3mVoice.vol0Opt")
	//voiceFilter
	//pluginFilter
}
