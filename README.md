# playback

![Tests](https://github.com/gotracker/playback/workflows/Tests/badge.svg)

## What is it?

It's an embeddable tracked music player written in Go.

## Why does this exist?

[Heucuva](https://github.com/heucuva/) needed to learn Go forever ago and figured this was a good way to do it. Also, the [Gotracker](https://github.com/gotracker/gotracker) project started growing into something more than just a command line player, so the rendering portion got ripped out and moved here.

## What does it play?

Files from/of the following formats/trackers:
* S3M - ScreamTracker 3
* MOD - Protracker/Fasttracker/Startrekker (_internally up-converted to S3M_)
* XM - Fasttracker II
* IT - Impulse Tracker

## What systems does it work on?

* Any, so long as you provide a way for the playback system to output its rendered content to somewhere useful.

## Requirements
* Go v1.18 or newer

## How does it work?

Not well, but it's good enough to play some moderately complex stuff.

## How do I use it?

Take a look at a few examples provided in the [internal/examples](internal/examples) folder. They will be able to show a start-to-finish example of the player, final stage mixing, and format conversion code in action.

## Bugs

### Known bugs

| Tags | Notes |
|------|-------|
| `player` | Unknown/unhandled commands (effects) will cause a panic. There aren't many left, but there are still some laying around. |
| `player` | The rendering system is fairly bad - it originally was designed only to work with S3M, but we decided to rework some of it to be more flexible. We managed to pull most of the mixing functionality out into somewhat generic structures/algorithms, but it still needs a lot of work. |
| `loader` | Attempting to load a corrupted tracker file may cause the deserializer to panic or go running off into the weeds indefinitely. |
| `mod` | MOD file support is buggy, at best. |
| `mod` `loader` | MOD files are up-converted to S3M internally and the S3M player uses NTSC-based lookup tables, so with a PAL-based MOD, the period values produced will end up being very slightly divergent from what is expected, as the S3M format converts note information to key-octave pairs, opting to look up the period information at time of need instead. |
| `xm` | XM file support is in a somewhat nascent state. Playback should work alright, but some things like Linear Frequency Slides are a little rough. |
| `it` | IT file support is in a somewhat nascent state. Playback should work alright in most cases, but some things like DSP plugins will not function. |
| `s3m` `opl2` | Attempting to play an S3M file with Adlib/OPL2 instruments does not produce the expected output. The OPL2 code has something wrong with it - it sounds pretty bad, though steps have been taken to remedy its strange output. |
| `mod` `s3m` | Amiga Paula/"LED" low-pass filter support is available, but the filter itself is a very lazy (and very over-optimized) Butterworth implementation. It will not produce the expected output. |
| `s3m` | SoundBlaster low-pass filter support is available, but comes in the form of a reused Amiga Paula low-pass (3.2kHz) filter. It does not function on the final output data, but instead the separate pre-final output channels. Taking all that into account, the output will not match expectations, but will perform relatively ok. |
| `xm` `opl2` | Attempting to play an XM file with Adlib/OPL2 instruments does not work. Most of the code for playback is there, but there's none for loading OPL2 instruments from file, so there's no way for the instruments to make it to the playback code. |
| `player` | Channel readouts are lazily attempted to match the layout from the tracker the song file came from. As a result, there are probably strange artifacts presented in it by the attempted simulation. |
| `player` `mixing` | The mixer still uses some simple saturation mixing techniques, but it's a lot better than it used to be. |
| `xm` `it` | Linear Frequency Slide support uses an _in-situ_ floating point power-of-2 calculation, which may be very slow on some hardware. Additionally, it is not going to match what Fasttracker II and Impulse Tracker do internally - using a pre-calculated lookup table - so the output may sound slightly different from expectation. |

### Unknown bugs

* There are many, we're sure.

## Further reading

Take a look at the fmoddoc2 documentation that the folks at FireLight studios released forever ago - it has great info how how to make a mod player, upgrade it to an s3m player, and then dork around with the internals a bit.
