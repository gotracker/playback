package channel

type SharedMemory struct {
	VolSlideEveryTick   bool
	ST300Portas         bool
	LowPassFilterEnable bool
	// ResetMemoryAtStartOfOrder0 if true will reset the memory registers when the first tick of the first row of the first order pattern plays
	ResetMemoryAtStartOfOrder0 bool
	// ST2/Amiga quirks
	ST2Vibrato          bool
	ST2Tempo            bool
	AmigaSlides         bool
	ZeroVolOptimization bool
	AmigaLimits         bool
	// Mod quirks mode
	ModCompatibility bool
}
