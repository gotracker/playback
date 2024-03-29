package loop

// Disabled is a disabled loop
//  |start>----------------------------------------end| <= on playthrough 1, whole sample plays
type Disabled struct{}

// Enabled returns true if the loop is enabled
func (l *Disabled) Enabled() bool {
	return false
}

// Length returns the length of the loop
func (l *Disabled) Length() int {
	return 0
}

// CalcPos calculates the position based on the loop details
func (l *Disabled) CalcPos(pos int, length int) (int, bool) {
	return min(max(pos, 0), length), false
}
