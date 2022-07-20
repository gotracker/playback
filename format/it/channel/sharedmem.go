package channel

type SharedMemory struct {
	// LinearFreqSlides is true if linear frequency slides are enabled (false = amiga-style period-based slides)
	LinearFreqSlides bool
	// OldEffectMode performs somewhat different operations for some effects:
	// On:
	//  - Vibrato does not operate on tick 0 and has double depth
	//  - Sample Offset will ignore the command if it would exceed the length
	// Off:
	//  - Vibrato is updated every frame
	//  - Sample Offset will set the offset to the end of the sample if it would exceed the length
	OldEffectMode bool
	// EFGLinkMode will make effects Exx, Fxx, and Gxx share the same memory
	EFGLinkMode bool
	// ResetMemoryAtStartOfOrder0 if true will reset the memory registers when the first tick of the first row of the first order pattern plays
	ResetMemoryAtStartOfOrder0 bool
}
