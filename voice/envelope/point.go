package envelope

// Point is a point for the envelope
type Point[T any] struct {
	Pos    int
	Length int
	Y      T
}
