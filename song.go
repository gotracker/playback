package playback

import "github.com/gotracker/playback/song"

// Song is a song.Data that can be constructed into a player that can play it
type Song interface {
	song.Data
	ConstructPlayer() (Playback, error)
}
