// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

func flipDesignatorRegisterBit31(e *InstructionEngine) {
	dr := e.GetDesignatorRegister()
	dr.SetBasicModeBaseRegisterSelection(!dr.GetBasicModeBaseRegisterSelection())
}

// unconditional jumps -------------------------------------------------------------------------------------------------

// StoreLocationAndJump (SLJ) stores the relative address of the instruction incremented by one, in the location
// specified by U, then loads the program counter with U+1.
func StoreLocationAndJump(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	value := (e.GetProgramAddressRegister().GetProgramCounter() + 1) & 0_777777
	e.StoreOperand(false, true, false, false, value)

	var flip31 bool
	var operand uint64
	operand, flip31, completed, interrupt = e.GetJumpOperand()
	if completed && interrupt == nil {
		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, false) // we need auto-increment to get us to the next instruction
		if flip31 {
			flipDesignatorRegisterBit31(e)
		}
	}

	return
}

// LoadModifierAndJump (LMJ) stores the incremented-by-one of the instruction's relative address into
// the 18-bit modifier portion of Xa, then loads the program counter from the U field.
func LoadModifierAndJump(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	var flip31 bool
	var operand uint64
	operand, flip31, completed, interrupt = e.GetJumpOperand()

	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		xReg := e.GetExecOrUserXRegister(ci.GetA())
		xReg.SetXM(e.GetProgramAddressRegister().GetProgramCounter() + 1)

		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, true)
		if flip31 {
			flipDesignatorRegisterBit31(e)
		}
	}

	return
}

// Jump (J) Loads the program counter from the U field - assumes no bank switching
func Jump(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	var flip31 bool
	var operand uint64
	operand, flip31, completed, interrupt = e.GetJumpOperand()

	if completed && interrupt == nil {
		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, true)
		if flip31 {
			flipDesignatorRegisterBit31(e)
		}
	}

	return
}

// JumpKeys (JK) evaluates the operand, but does not jump.
// The assumption is that the selected jump key is present, but cleared.
// It is not specified how the jump key is selected, but it doesn't matter.
func JumpKeys(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	_, _, completed, interrupt = e.GetJumpOperand()
	return
}

// HaltJump (HLTJ) Loads the program counter for the U field then stops the processor
func HaltJump(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	if e.activityStatePacket.GetDesignatorRegister().GetProcessorPrivilege() > 0 {
		completed = false
		interrupt = pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
		return
	}

	var flip31 bool
	var operand uint64
	operand, flip31, completed, interrupt = e.GetJumpOperand()

	if completed && interrupt == nil {
		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, true)
		if flip31 {
			flipDesignatorRegisterBit31(e)
		}
		e.Stop(HaltJumpExecutedStop, 0)
	}

	return
}

// HaltKeysAndJump (HJ or HKJ) Loads the program counter for the U field. No halt occurs.
func HaltKeysAndJump(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	var flip31 bool
	var operand uint64
	operand, flip31, completed, interrupt = e.GetJumpOperand()

	if completed && interrupt == nil {
		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, true)
		if flip31 {
			flipDesignatorRegisterBit31(e)
		}
	}

	return
}

// Jumps conditional upon some value -----------------------------------------------------------------------------------

// JumpZero (JZ) Loads the program counter from the U field *IF* Aa is zero
func JumpZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsZero() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// DoubleJumpZero (DJZ) Loads the program counter from the U field *IF* Aa | Aa+1 is zero
func DoubleJumpZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	ci := e.GetCurrentInstruction()
	ax := ci.GetA()
	aReg1 := e.GetExecOrUserARegister(ax).GetW()
	aReg2 := e.GetExecOrUserARegister(ax + 1).GetW()

	if ((aReg1 == pkg.PositiveZero) && (aReg2 == pkg.PositiveZero)) ||
		((aReg1 == pkg.NegativeZero) && (aReg2 == pkg.NegativeZero)) {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpNonZero (JNZ) Loads the program counter from the U field *IF* Aa is not zero
func JumpNonZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if !aReg.IsZero() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpPositive (JP) Loads the program counter from the U field *IF* Aa is positive
func JumpPositive(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsPositive() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpPositiveAndShift (JPS) Loads the program counter from the U field *IF* Aa is positive.
// Aa is shifted left circularly by one bit in any case
func JumpPositiveAndShift(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsPositive() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if !completed || interrupt != nil {
			return
		}

		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, true)
		if flip31 {
			flipDesignatorRegisterBit31(e)
		}
	}

	aReg.ShiftLeftCircular(1)
	return
}

// JumpNegative (JN) Loads the program counter from the U field *IF* Aa is negative
func JumpNegative(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsNegative() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpNegativeAndShift (JNS) Loads the program counter from the U field *IF* Aa is negative.
// Aa is shifted left circularly by one bit in any case
func JumpNegativeAndShift(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsNegative() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if !completed || interrupt != nil {
			return
		}

		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, true)
		if flip31 {
			flipDesignatorRegisterBit31(e)
		}
	}

	aReg.ShiftLeftCircular(1)
	return
}

// JumpLowBit (JB) Loads the program counter from the U field *IF* the least significant bit of Aa is set
func JumpLowBit(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if uint64(*aReg)&01 != 0 {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpNoLowBit (JNB) Loads the program counter from the U field *IF* the least significant bit of Aa is clear
func JumpNoLowBit(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if uint64(*aReg)&01 == 0 {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpGreaterAndDecrement (JGD) Loads the program counter from the U field *IF* the GRS register indicated by
// the value created by the right-most seven bits of the concatenation of the j-field to the a-field is > 0.
// In any case, the indicated register is decremented after the comparison to zero.
func JumpGreaterAndDecrement(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	ci := e.GetCurrentInstruction()
	ix := ((ci.GetJ() << 4) | ci.GetA()) & 0177
	value := e.generalRegisterSet.registers[ix]
	if value.IsPositive() && !value.IsZero() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if !completed || interrupt != nil {
			return
		}

		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, true)
		if flip31 {
			flipDesignatorRegisterBit31(e)
		}
	}

	newValue := pkg.AddSimple(value.GetW(), pkg.NegativeOne)
	e.generalRegisterSet.registers[ix].SetW(newValue)
	return
}

// JumpModifierGreaterAndIncrement (JMGI) Loads the program counter from the U field *IF* the modifier portion of Xa
// is greater than zero. In any case, the signed increment of Xa is added to the signed modifier of Xa
// (after the comparison to zero).
func JumpModifierGreaterAndIncrement(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	ci := e.GetCurrentInstruction()
	xReg := e.GetExecOrUserXRegister(ci.GetA())
	reg := pkg.Word36(*xReg)
	modifier := reg.GetXH2()
	if pkg.IsPositive(modifier) && !pkg.IsZero(modifier) {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if !completed || interrupt != nil {
			return
		}

		e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
		e.SetProgramCounter(operand, true)
		if flip31 {
			flipDesignatorRegisterBit31(e)
		}
	}

	xReg.IncrementModifier()
	return
}

// Jumps conditional upon a designator register bit --------------------------------------------------------------------

// JumpOverflow (JO) Loads the program counter from the U field *IF* designator register bit 19 is set.
func JumpOverflow(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	if e.activityStatePacket.GetDesignatorRegister().IsOverflowSet() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpNoOverflow (JNO) Loads the program counter from the U field *IF* designator register bit 19 is clear.
func JumpNoOverflow(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	if !e.activityStatePacket.GetDesignatorRegister().IsOverflowSet() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpCarry (JC) Loads the program counter from the U field *IF* designator register bit 18 is set.
func JumpCarry(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	if e.activityStatePacket.GetDesignatorRegister().IsCarrySet() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpNoCarry (JNC) Loads the program counter from the U field *IF* designator register bit 18 is clear.
func JumpNoCarry(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	if !e.activityStatePacket.GetDesignatorRegister().IsCarrySet() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpDivideFault (JDF) Loads the program counter from the U field *IF* designator register bit 23 is set.
func JumpDivideFault(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	if e.activityStatePacket.GetDesignatorRegister().IsDivideCheckSet() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpNoDivideFault (JNDF) Loads the program counter from the U field *IF* designator register bit 23 is clear.
func JumpNoDivideFault(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	if !e.activityStatePacket.GetDesignatorRegister().IsDivideCheckSet() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpFloatingOverflow (JFO) Loads the program counter from the U field *IF* designator register bit 22 is set.
func JumpFloatingOverflow(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	if e.activityStatePacket.GetDesignatorRegister().IsCharacteristicOverflowSet() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpNoFloatingOverflow (JNFO) Loads the program counter from the U field *IF* designator register bit 22 is clear.
func JumpNoFloatingOverflow(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	if !e.activityStatePacket.GetDesignatorRegister().IsCharacteristicOverflowSet() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpFloatingUnderflow (JFU) Loads the program counter from the U field *IF* designator register bit 21 is set.
func JumpFloatingUnderflow(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	if e.activityStatePacket.GetDesignatorRegister().IsCharacteristicUnderflowSet() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}

// JumpNoFloatingUnderflow (JNFU) Loads the program counter from the U field *IF* designator register bit 21 is clear.
func JumpNoFloatingUnderflow(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	if !e.activityStatePacket.GetDesignatorRegister().IsCharacteristicUnderflowSet() {
		var flip31 bool
		var operand uint64
		operand, flip31, completed, interrupt = e.GetJumpOperand()

		if completed && interrupt == nil {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}
	}

	return
}
