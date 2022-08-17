package feature

type MovingAverageFilter struct {
	Enabled bool
	// WindowSize cannot be 0!
	WindowSize int
}
