package filter

import (
	"math"

	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/frequency"
)

type amigaLPFChannelData struct {
	ynz1 volume.Volume
	ynz2 volume.Volume
}

// AmigaLPF is a 12dB/octave 2-pole Butterworth Low-Pass Filter with 3275 Hz cut-off
type AmigaLPF struct {
	channels []amigaLPFChannelData
	a0       volume.Volume
	b0       volume.Volume
	b1       volume.Volume

	playbackRate frequency.Frequency
}

// NewAmigaLPF creates a new AmigaLPF
func NewAmigaLPF(instrument frequency.Frequency) *AmigaLPF {
	var f AmigaLPF
	f.SetPlaybackRate(instrument)
	return &f
}

func (f *AmigaLPF) Clone() Filter {
	c := *f
	c.channels = make([]amigaLPFChannelData, len(f.channels))
	for i := range f.channels {
		c.channels[i] = f.channels[i]
	}
	return &c
}

func (f *AmigaLPF) SetPlaybackRate(playback frequency.Frequency) {
	if f.playbackRate == playback {
		return
	}
	f.playbackRate = playback

	freq := 3275.0

	f2 := float64(playback) / 2.0
	if freq > f2 {
		freq = f2
	}

	fc := freq * 2.0 * math.Pi

	r := float64(playback) / fc

	d := r
	e := r * r

	a := 1.0 / (1.0 + d + e)
	b := (d + e + e) * a
	c := -e * a

	f.a0 = volume.Volume(a)
	f.b0 = volume.Volume(b)
	f.b1 = volume.Volume(c)
}

// Filter processes incoming (dry) samples and produces an outgoing filtered (wet) result
func (f *AmigaLPF) Filter(dry volume.Matrix) volume.Matrix {
	if dry.Channels == 0 {
		return volume.Matrix{}
	}
	wet := dry // we can update in-situ and be ok
	for i := 0; i < dry.Channels; i++ {
		s := dry.StaticMatrix[i]
		for len(f.channels) <= i {
			f.channels = append(f.channels, amigaLPFChannelData{})
		}
		c := &f.channels[i]

		xn := s
		yn := min(max(xn*f.a0+c.ynz1*f.b0+c.ynz2*f.b1, -1), 1)
		c.ynz2 = c.ynz1
		c.ynz1 = yn
		wet.StaticMatrix[i] = yn
	}
	return wet
}

// UpdateEnv updates the filter with the value from the filter envelope
func (f *AmigaLPF) UpdateEnv(v uint8) {
}
