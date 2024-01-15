package opl2

import (
	"github.com/gotracker/opl2"
)

// Chip sets up a contract that the chip definition will contain these interfaces
type Chip interface {
	WriteReg(uint32, uint8)
	GenerateBlock2(uint, []int32)
}

// OPL2Intf is an interface to get the active OPL2 chip
type OPL2Intf interface {
	GetOPL2Chip() Chip
}

// NewOPL2Chip generates a new OPL2 chip
func NewOPL2Chip(rate uint32) Chip {
	return opl2.NewChip(rate, false)
}
