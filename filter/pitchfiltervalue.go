package filter

type PitchFiltValue uint8

func (p PitchFiltValue) AsFilter() uint8 {
	return uint8(p)
}

func (p PitchFiltValue) AsPitch() int8 {
	return int8(p)
}
