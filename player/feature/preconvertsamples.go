package feature

import "github.com/gotracker/playback/voice/pcm"

type PreConvertSamples struct {
	Enabled       bool
	DesiredFormat pcm.SampleDataFormat
}

func UseNativeSampleFormat(enabled bool) Feature {
	return PreConvertSamples{
		Enabled:       enabled,
		DesiredFormat: pcm.SampleDataFormatNative,
	}
}
