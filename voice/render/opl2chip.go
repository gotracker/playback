package render

// OPL2Chip sets up a contract that the chip definition will contain these interfaces
type OPL2Chip interface {
	WriteReg(uint32, uint8)
	GenerateBlock2(uint, []int32)
}

// OPL2Intf is an interface to get the active OPL2 chip
type OPL2Intf interface {
	GetOPL2Chip() OPL2Chip
}
