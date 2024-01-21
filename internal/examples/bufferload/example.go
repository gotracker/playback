package main

import (
	"bytes"
	"errors"
	"os"

	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/playback/format"
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/sampler"
	"github.com/gotracker/playback/song"
)

// ExamplePlayBufferToStdout will read in a music module file (in this case: `test/ode_to_protracker.mod`)
// and render it out in 44100Hz, 2 channel, 16-bit signed, little-endian integer PCM format.
// If you have PulseAudio installed and configured correctly, then simply pipe the output from this test
// into it with `go run ./internal/examples/fileload | pacat -p --channels=2 --rate=44100 --format=s16le`
func ExamplePlayBufferToStdout() {
	const (
		sampleRate   = 44100
		channels     = 2
		sampleFormat = sampling.Format16BitLESigned
	)

	// This is a list of features we can build up before handing off to the loader and player.
	var features []feature.Feature

	// Let's start off by adding some loader features.

	// In this case, we are going to enable native sample type conversion - this allows the loader to
	// automatically convert the sample data into a format that doesn't need ad-hoc conversion to a form
	// that the player can directly use. This conversion comes at the cost of large amounts of memory,
	// so be prepared for that. This feature is built upon another feature type, so it needs a helper
	// function to build it out.
	features = append(features, feature.UseNativeSampleFormat(true))

	// Let's assume we don't know who made the song we're trying to play and that the author used some
	// truly insane channel effects - under most circumstances, this would cause the player to panic,
	// but we don't want that for this test. We can squelch the panic and continue as if the strange
	// effect wasn't even there using this feature.
	features = append(features, feature.IgnoreUnknownEffect{Enabled: true})

	// Let's also say that we know there are built-in pattern loop commands in the song we want to play
	// (foreshadowing: there are) and that we want to only play the whole song once, thus stopping
	// playback when the first pattern loop command tries to go back to a part of the song we've already
	// played. This feature can prevent or allow that - even up to a set number of loop attempts.
	// In this case, we are going to set the loop count to 0, which disables the looping and only plays
	// the song one time through.
	features = append(features, feature.SongLoop{Count: 0})

	// There's an automagical loader utility which divines the file type and presents a song that can
	// effectively represent it. See `ExamplePlayFileToStdout` (./internal/examples/fileload) for an
	// example of the file loader version of this call. In this example, we know the format, so we can
	// pass in the specific format loader we want to use.
	songData, songFormat, err := format.LoadFromReader("mod", bytes.NewReader(modfile), features)
	if err != nil {
		panic(err)
	}

	// Here is where we get a final chance to submit any overrides or configurations we want to
	// supply - we can send it the configuration we already have built up, since it will know how to
	// pull the settings it wants, so no need to worry about filtering or splitting out the settings.
	var userSettings settings.UserSettings
	if err := songFormat.ConvertFeaturesToSettings(&userSettings, features); err != nil {
		panic(err)
	}

	// Next, create a player machine to operate over the song configured by the settings.
	player, err := machine.NewMachine(songData, userSettings)
	if err != nil {
		panic(err)
	}

	// Now that the player is configured, we can allocate a channel is filled with bundles of
	// pre-mix data (i.e.: data which is ready to be converted into the final, mixed version).
	// We want to use a buffered channel with a bit of room in it. It doesn't have to be huge,
	// but the more space you supply, the more likely the player will not underflow the channel.
	const premixChannelSize = 8
	premixDataChannel := make(chan *output.PremixData, premixChannelSize)

	// We need to set up a consumer for the pre-mix data. The goal is to output it to the standard
	// output device. Since we must use a goroutine to do this work, we need some sort of signalling
	// mechanism to handle when the song ends (or an error occurs).  Normally, we would use a context
	// but in this case, when the function ends, the channel deferred close will cause the goroutine
	// to end.
	defer close(premixDataChannel)

	// Now that we have a player allocated for the format, we need to tell it the minimal configuration
	// for the stream of data we are wanting to produce - namely, the sampling rate and the number of
	// channels. These first two parameters are fundamental to a huge number of operations, so they must
	// be set outside of the configuration process you will see below. The third parameter provides a
	// way for the calling application (our example) to get the generated output data in the form of
	// pre-mixed packets. These packets can be further mixed into audio streams for use with sound
	// devices and files.
	out := sampler.NewSampler(sampleRate, channels, func(premix *output.PremixData) {
		// put our premixed data into the premixDataChannel we built earlier.
		premixDataChannel <- premix
	})
	if out == nil {
		panic(errors.New("could not create sampler"))
	}

	// Our desire is to output a specific format of PCM audio data to the standard output device, so
	// we need to mix and convert the pre-mix data into that format. This mixer will be able to do
	// just that.
	m := mixing.Mixer{
		Channels: channels,
	}

	// A panning mixer knows how to coordinate panning values into a pre-final (penultimate?) mixing
	// matrix that can be collapsed into the final form, ready for the sample type conversion.
	panMixer := mixing.GetPanMixer(channels)

	go func() {
		// Wait for a pre-mix data blob to show up
		for premix := range premixDataChannel {
			// Flatten the data into the final format - this is a very complex process that this one
			// helper function miraculously does for us, placing into a very handy slice of bytes.
			data := m.Flatten(panMixer, premix.SamplesLen, premix.Data, premix.MixerVolume, sampleFormat)

			// write it out!  If we run into an error, then ignore it for now. This is where a context
			// with a cancellation would be a good solution to properly coordinate the player update
			// process (see below) with any critical errors we receive here.
			_, _ = os.Stdout.Write(data)
		}
	}()

playerUpdateLoop:
	for {
		// Now we need to tell the player to update its internal state - this will generate a single
		// row tick's worth of pre-mix data and call our callback function specified in the Sampler
		// stage we specified earlier. Normally, we would want to set up a goroutine for this call to
		// run in, but in this example, we're fine to do a simple loop.
		if err := player.Tick(out); err != nil {
			// In the event we finish our song, we will receive a specific error message informing us
			// we can quit.
			if errors.Is(err, song.ErrStopSong) {
				break playerUpdateLoop
			}
			// If we get here, then we don't know what the error is...
			panic(err)
		}
	}

	// We're done!
}

func main() {
	ExamplePlayBufferToStdout()
}
