// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

// StoreLocationAndJump (SLJ) stores the relative address of the instruction incremented by one in the location
// specified by U, then loads the program counter with U+1.
func StoreLocationAndJump(e *InstructionEngine) (completed bool) {
	value := (e.GetProgramAddressRegister().GetProgramCounter() + 1) & 0_777777
	comp1, i1 := e.StoreOperand(false, true, false, false, value)
	if i1 != nil {
		e.PostInterrupt(i1)
		return false
	} else if !comp1 {
		return false
	}

	operand, flip31, comp2, i2 := e.GetJumpOperand()
	if i2 != nil {
		e.PostInterrupt(i2)
		return false
	} else if !comp2 {
		return false
	}

	e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
	e.SetProgramCounter(operand, false) // we need auto-increment to get us to the next instruction
	if flip31 {
		e.flipDesignatorRegisterBit31()
	}

	return true
}

// LoadModifierAndJump (LMJ) stores the incremented-by-one of the instruction's relative address into
// the 18-bit modifier portion of Xa, then loads the program counter from the U field.
func LoadModifierAndJump(e *InstructionEngine) (completed bool) {
	operand, flip31, comp, i := e.GetJumpOperand()

	if i != nil {
		e.PostInterrupt(i)
	} else if comp {
		ci := e.GetCurrentInstruction()
		xReg := e.GetExecOrUserXRegister(ci.GetA())
		xReg.SetXM(e.GetProgramAddressRegister().GetProgramCounter() + 1)

		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, true)
		if flip31 {
			e.flipDesignatorRegisterBit31()
		}
	}

	return comp
}

// Jump (J) Loads the program counter from the U field - assumes no bank switching
func Jump(e *InstructionEngine) (completed bool) {
	operand, flip31, comp, i := e.GetJumpOperand()

	if i != nil {
		e.PostInterrupt(i)
	} else if comp {
		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, true)
		if flip31 {
			e.flipDesignatorRegisterBit31()
		}
	}

	return comp
}

// JumpKeys (JK) evaluates the operand, but does not jump.
// The assumption is that the selected jump key is present, but cleared.
// It is not specified how the jump key is selected, but it doesn't matter.
func JumpKeys(e *InstructionEngine) (completed bool) {
	_, _, comp, i := e.GetJumpOperand()

	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// HaltJump (HLTJ) Loads the program counter for the U field then stops the processor
func HaltJump(e *InstructionEngine) (completed bool) {
	if e.activityStatePacket.GetDesignatorRegister().GetProcessorPrivilege() > 0 {
		i := pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
		e.PostInterrupt(i)
		return false
	}

	operand, flip31, comp, i := e.GetJumpOperand()

	if i != nil {
		e.PostInterrupt(i)
	} else if comp {
		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, true)
		if flip31 {
			e.flipDesignatorRegisterBit31()
		}
		e.Stop(HaltJumpExecutedStop, 0)
	}

	return comp
}

// HaltKeysAndJump (HJ or HKJ) Loads the program counter for the U field. No halt occurs.
func HaltKeysAndJump(e *InstructionEngine) (completed bool) {
	operand, flip31, comp, i := e.GetJumpOperand()

	if i != nil {
		e.PostInterrupt(i)
	} else if comp {
		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, true)
		if flip31 {
			e.flipDesignatorRegisterBit31()
		}
	}

	return comp
}
