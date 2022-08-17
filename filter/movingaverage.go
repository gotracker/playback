package filter

import (
	"math"

	"github.com/gotracker/gomixing/volume"
)

type movavgPoint struct {
	data  volume.Matrix
	coeff volume.Volume
}

type MovingAverage struct {
	points []movavgPoint
}

func NewMovingAverage(windowSize int) Filter {
	if windowSize == 0 {
		panic("windowSize cannot be 0")
	}
	ma := MovingAverage{
		points: make([]movavgPoint, windowSize),
	}

	sigma := float64(windowSize)
	a := 1.0 / (sigma * math.Sqrt(2.0*math.Pi))
	mu := float64(windowSize) / 2.0
	sigmasq2 := 2 * sigma * sigma
	for x := 0; x < windowSize; x++ {
		xmu := float64(x) - mu
		fx := a * math.Exp(-(xmu*xmu)/sigmasq2)
		if fx < 0 {
			fx = 0
		}
		if fx > 1 {
			fx = 1
		}
		ma.points[x].coeff = volume.Volume(fx)
	}
	return &ma
}

func (ma *MovingAverage) Clone() Filter {
	clone := MovingAverage{
		points: make([]movavgPoint, len(ma.points)),
	}
	copy(clone.points, ma.points)
	return &clone
}

func (ma *MovingAverage) Filter(dry volume.Matrix) volume.Matrix {
	var wet volume.Matrix
	// shift the data down by 1
	windowLen := len(ma.points)
	var lastT volume.Volume
	for i := 1; i < windowLen; i++ {
		prevPt := &ma.points[i-1]
		pt := &ma.points[i]
		prevPt.data = pt.data
		wet.Accumulate(pt.data.Apply(prevPt.coeff))
		lastT = pt.coeff
	}
	// now set our dry data into the buffer
	ma.points[windowLen-1].data = dry
	// copy our dry data into our wet matrix
	wet.Accumulate(dry.Apply(lastT))

	return wet
}

func (ma *MovingAverage) UpdateEnv(val int8) {

}
