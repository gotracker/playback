package instrument

import (
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/voice/pcm"
)

func NewSample(data []byte, length int, channels int, format pcm.SampleDataFormat, features []feature.Feature) (pcm.Sample, error) {
	sf := format
	for _, feat := range features {
		switch f := feat.(type) {
		case feature.PreConvertSamples:
			if f.Enabled {
				sf = f.DesiredFormat
			}
		}
	}

	if sf == format {
		// original format
		return pcm.NewSample(data, length, channels, format), nil
	}

	inSample := pcm.NewSample(data, length, channels, format)
	if sf != pcm.SampleDataFormatNative {
		// format conversion
		outSample, err := pcm.ConvertTo(inSample, sf)
		if err != nil {
			return nil, err
		}
		return outSample, nil
	}

	// native conversion
	nativeData := make([]volume.Matrix, 0, length)
	for i := 0; i < length; i++ {
		d, err := inSample.Read()
		if err != nil {
			return nil, err
		}
		nativeData = append(nativeData, d)
	}
	return pcm.NewSampleNative(nativeData, length, channels), nil
}
