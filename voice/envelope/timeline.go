package envelope

import (
	"container/heap"
	"math"

	"github.com/heucuva/optional"
)

type Timeline[T any] struct {
	p pointHeap[T]
}

func (t *Timeline[T]) Init() {
	heap.Init(&t.p)
}

func (t *Timeline[T]) Push(pos int, value T) {
	heap.Push(&t.p, point[T]{
		AddPos: t.p.Len(),
		Pos:    pos,
		Y:      value,
	})
}

func (t Timeline[T]) Result() []EnvPoint[T] {
	var (
		result    []EnvPoint[T]
		prevPoint optional.Value[point[T]]
	)
	length := t.p.Len()
	for i := 0; i < length; i++ {
		entry := heap.Pop(&t.p).(point[T])
		if prev, prevSet := prevPoint.Get(); prevSet {
			next := EnvPoint[T]{
				Ticks: entry.Pos - prev.Pos,
				Y:     prev.Y,
			}
			result = append(result, next)
		}
		prevPoint.Set(entry)
	}
	var last point[T]
	last, _ = prevPoint.Get()
	result = append(result, EnvPoint[T]{
		Ticks: math.MaxInt,
		Y:     last.Y,
	})
	return result
}

type pointHeap[T any] []point[T]

func (p pointHeap[T]) Len() int {
	return len(p)
}

func (p pointHeap[T]) Less(i, j int) bool {
	lhs := p[i]
	rhs := p[j]
	if lhs.Pos > rhs.Pos {
		return false
	} else if lhs.Pos == rhs.Pos {
		return lhs.AddPos < rhs.AddPos
	}
	return true
}

func (p *pointHeap[T]) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *pointHeap[T]) Push(x any) {
	*p = append(*p, x.(point[T]))
}

func (p *pointHeap[T]) Pop() any {
	lpos := p.Len()
	last := (*p)[lpos-1]
	*p = (*p)[:lpos-1]
	return last
}

type point[T any] struct {
	AddPos int
	Pos    int
	Y      T
}
