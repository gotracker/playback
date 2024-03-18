package mixing

import (
	"bytes"
	"time"

	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
)

// SampleMixIn is the parameters for mixing in a sample into a MixBuffer
type SampleMixIn struct {
	Sample    sampling.Sampler
	StaticVol volume.Volume
	PanMatrix panning.PanMixer
	MixPos    int
	MixLen    int
}

// MixBuffer is a buffer of premixed volume data intended to
// be eventually sent out to the sound output device after
// conversion to the output format
type MixBuffer []volume.Matrix

// C returns a channel and a function that flushes any outstanding mix-ins and closes the channel
func (m *MixBuffer) C() (chan<- SampleMixIn, func()) {
	ch := make(chan SampleMixIn, 32)
	go func() {
		for d := range ch {
			m.MixInSample(d)
		}
	}()
	return ch, func() {
		for len(ch) != 0 {
			time.Sleep(1 * time.Millisecond)
		}
		close(ch)
	}
}

// MixInSample mixes in a single sample entry into the mix buffer
func (m *MixBuffer) MixInSample(d SampleMixIn) {
	pos := d.MixPos
	for i := 0; i < d.MixLen; i++ {
		dry := d.Sample.GetSample()
		samp := dry.Apply(d.StaticVol)
		mixed := d.PanMatrix.ApplyToMatrix(samp)
		(*m)[pos].Accumulate(mixed)
		pos++
		d.Sample.Advance()
	}
}

// Add will mix in another MixBuffer's data
func (m *MixBuffer) Add(pos int, rhs MixBuffer, volMtx volume.Matrix) {
	maxLen := len(rhs)
	for i := 0; i < maxLen; i++ {
		out := volMtx.ApplyToMatrix(rhs[i])
		(*m)[pos+i].Accumulate(out)
	}
}

// ToRenderData converts a mixbuffer into a byte stream intended to be
// output to the output sound device
func (m *MixBuffer) ToRenderData(samples int, channels int, mixerVolume volume.Volume, formatter sampling.Formatter) []byte {
	writer := &bytes.Buffer{}
	writer.Grow(samples * ((formatter.Size() + 7) / 8) * channels)
	for _, samp := range *m {
		buf := samp.Apply(mixerVolume)
		d := buf.ToChannels(channels)
		for i := 0; i < channels; i++ {
			_ = formatter.Write(writer, d.StaticMatrix[i]) // lint
		}
	}
	return writer.Bytes()
}

// ToIntStream converts a mixbuffer into an int stream intended to be
// output to the output sound device
func (m *MixBuffer) ToIntStream(outputChannels int, samples int, bitsPerSample int, mixerVolume volume.Volume) [][]int32 {
	data := make([][]int32, outputChannels)
	for c := range data {
		data[c] = make([]int32, samples)
	}
	for i, samp := range *m {
		buf := samp.Apply(mixerVolume)
		d := buf.ToChannels(outputChannels)
		for c := 0; c < outputChannels; c++ {
			data[c][i] = d.StaticMatrix[c].ToIntSample(bitsPerSample)
		}
	}
	return data
}

// ToRenderDataWithBufs converts a mixbuffer into a byte stream intended to be
// output to the output sound device
func (m *MixBuffer) ToRenderDataWithBufs(outBuffers [][]byte, samples int, mixerVolume volume.Volume, formatter sampling.Formatter) {
	pos := 0
	onum := 0
	out := outBuffers[onum]
	for _, samp := range *m {
		buf := samp.Apply(mixerVolume)
		for c := 0; c < buf.Channels; c++ {
			for pos >= len(out) {
				onum++
				if onum > len(outBuffers) {
					return
				}
				out = outBuffers[onum]
				pos = 0
			}
			_ = formatter.WriteAt(out, int64(pos), buf.StaticMatrix[c]) // lint
			pos += formatter.Size()
		}
	}
}
