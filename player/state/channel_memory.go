package state

type ChannelMemory[TMemory any] struct {
	Memory *TMemory
}

// GetMemory returns the interface to the custom effect memory module
func (cm *ChannelMemory[TMemory]) GetMemory() *TMemory {
	return cm.Memory
}

// SetMemory sets the custom effect memory interface
func (cm *ChannelMemory[TMemory]) SetMemory(mem *TMemory) {
	cm.Memory = mem
}
