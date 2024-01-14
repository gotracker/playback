package filter

import (
	"math"

	"github.com/gotracker/playback/period"

	"github.com/gotracker/gomixing/volume"
)

type EchoFilterSettings struct {
	WetDryMix  float32
	Feedback   float32
	LeftDelay  float32
	RightDelay float32
	PanDelay   float32
}

type EchoFilterFactory struct {
	Reserved00 [4]byte
	EchoFilterSettings
}

func (e *EchoFilterFactory) Factory() Factory {
	return func(instrument period.Frequency) Filter {
		echo := EchoFilter{
			EchoFilterSettings: e.EchoFilterSettings,
		}
		return &echo
	}
}

type delayInfo struct {
	buf   []volume.Volume
	delay int
}

//===========

type EchoFilter struct {
	EchoFilterSettings
	initialFeedback volume.Volume
	writePos        int
	delay           [2]delayInfo // L,R
}

func (e *EchoFilter) SetPlaybackRate(playback period.Frequency) {
	e.initialFeedback = volume.Volume(math.Sqrt(float64(1.0 - (e.Feedback * e.Feedback))))

	playbackRate := float32(playback)
	bufferSize := int(playbackRate * 2)

	for c, delayMs := range [2]float32{e.LeftDelay, e.RightDelay} {
		delay := int(delayMs * 2.0 * playbackRate)
		e.delay[c].delay = delay
		e.delay[c].buf = make([]volume.Volume, bufferSize)
	}
}

func (e *EchoFilter) Clone() Filter {
	clone := EchoFilter{
		EchoFilterSettings: e.EchoFilterSettings,
		writePos:           e.writePos,
	}
	for i := range clone.delay {
		clone.delay[i].buf = make([]volume.Volume, len(e.delay[i].buf))
		copy(clone.delay[i].buf, e.delay[i].buf)
	}
	return &clone
}

func (e *EchoFilter) Filter(dry volume.Matrix) volume.Matrix {
	if dry.Channels == 0 {
		return volume.Matrix{}
	}
	wetMix := volume.Volume(e.WetDryMix)
	dryMix := 1 - wetMix
	wet := dry.Apply(dryMix)

	feedback := volume.Volume(e.Feedback)

	crossEcho := e.PanDelay >= 0.5

	bufferLen := len(e.delay[0].buf)

	for e.writePos >= bufferLen {
		e.writePos -= bufferLen
	}
	for e.writePos < 0 {
		e.writePos += bufferLen
	}

	for c := 0; c < dry.Channels; c++ {
		readChannel := c
		if crossEcho {
			readChannel = 1 - c
		}
		read := &e.delay[readChannel]
		write := &e.delay[c]

		readPos := e.writePos - read.delay
		for readPos < 0 {
			readPos += bufferLen
		}
		for readPos >= bufferLen {
			readPos -= bufferLen
		}

		chnInput := dry.StaticMatrix[c]
		chnDelay := read.buf[readPos]

		chnOutput := chnInput * e.initialFeedback
		chnOutput += chnDelay * feedback

		write.buf[e.writePos] = chnOutput

		wet.StaticMatrix[c] += chnDelay * wetMix
	}

	e.writePos++

	return wet
}

func (e *EchoFilter) UpdateEnv(val uint8) {

}
