package util

import "golang.org/x/exp/constraints"

type Lerpable interface {
	constraints.Integer | constraints.Float
}

func Lerp[T Lerpable](t float64, a, b T) T {
	if t <= 0 {
		return a
	} else if t >= 1 {
		return b
	}
	return a + T(t*(float64(b)-(float64(a))))
}
