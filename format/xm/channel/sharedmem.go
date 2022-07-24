package channel

type SharedMemory struct {
	// LinearFreqSlides is true if linear frequency slides are enabled (false = amiga-style period-based slides)
	LinearFreqSlides bool
	// ExtendedFilterRange is true if the extended filter range is enabled
	ExtendedFilterRange bool
	// ResetMemoryAtStartOfOrder0 if true will reset the memory registers when the first tick of the first row of the first order pattern plays
	ResetMemoryAtStartOfOrder0 bool
}
