package machine

import (
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/player/machine/instruction"
)

type instructionRowStart[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] interface {
	instruction.Instruction
	RowStart(ch index.Channel, m Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoInstructionRowStart(ch index.Channel, i instruction.Instruction) error {
	ii, ok := i.(instructionRowStart[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning])
	if !ok {
		return nil
	}
	return ii.RowStart(ch, m)
}

type instructionPreTick[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] interface {
	instruction.Instruction
	PreTick(ch index.Channel, m Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], tick int) error
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoInstructionPreTick(ch index.Channel, i instruction.Instruction) error {
	ii, ok := i.(instructionPreTick[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning])
	if !ok {
		return nil
	}
	return ii.PreTick(ch, m, m.ticker.current.tick)
}

type instructionTick[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] interface {
	instruction.Instruction
	Tick(ch index.Channel, m Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], tick int) error
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoInstructionTick(ch index.Channel, i instruction.Instruction) error {
	ii, ok := i.(instructionTick[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning])
	if !ok {
		return nil
	}
	return ii.Tick(ch, m, m.ticker.current.tick)
}

type instructionPostTick[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] interface {
	instruction.Instruction
	PostTick(ch index.Channel, m Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], tick int) error
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoInstructionPostTick(ch index.Channel, i instruction.Instruction) error {
	ii, ok := i.(instructionPostTick[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning])
	if !ok {
		return nil
	}
	return ii.PostTick(ch, m, m.ticker.current.tick)
}

type instructionRowEnd[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] interface {
	instruction.Instruction
	RowEnd(ch index.Channel, m Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoInstructionRowEnd(ch index.Channel, i instruction.Instruction) error {
	ii, ok := i.(instructionRowEnd[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning])
	if !ok {
		return nil
	}
	return ii.RowEnd(ch, m)
}
