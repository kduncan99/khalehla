// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

// StoreLocationAndJump (SLJ) stores the relative address of the instruction incremented by one, in the location
// specified by U, then loads the program counter with U+1.
func StoreLocationAndJump(e *InstructionEngine) (bool, pkg.Interrupt) {
	value := e.relativeAddress + 1
	if value > 0_777777 {
		value = 0
	}

	e.StoreOperand(false, true, false, false, value)
	_, operand, interrupt := e.GetJumpOperand(true)
	if interrupt != nil {
		return false, interrupt
	}

	e.SetProgramCounter(operand, false) // we need auto-increment to get us to the next instruction
	return true, nil
}

// LoadModifierAndJump (LMJ) stores the incremented-by-one of the instruction's relative address into
// the 18-bit modifier portion of Xa, then loads the program counter from the U field.
func LoadModifierAndJump(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	xReg := e.GetExecOrUserXRegister(ci.GetA())
	value := e.relativeAddress + 1
	if value > 0_777777 {
		value = 0
	}
	xReg.SetXM(value)

	e.SetProgramCounter(operand, true)
	return true, nil
}

// Jump (J) Loads the program counter from the U field - assumes no bank switching
func Jump(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	e.SetProgramCounter(operand, true)
	return true, nil
}

// JumpKeys (JK) evaluates the operand, but does not jump.
// The assumption is that the selected jump key is present, but cleared.
// It is not specified how the jump key is selected, but it doesn't matter.
func JumpKeys(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, _, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	return true, nil
}

// HaltJump (HLTJ) Loads the program counter for the U field then stops the processor
func HaltJump(e *InstructionEngine) (bool, pkg.Interrupt) {
	if e.activityStatePacket.GetDesignatorRegister().GetProcessorPrivilege() > 0 {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
	}

	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	e.SetProgramCounter(operand, true)
	e.Stop(HaltJumpExecutedStop, 0)
	return true, nil
}

// HaltKeysAndJump (HJ or HKJ) Loads the program counter for the U field. No halt occurs.
func HaltKeysAndJump(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	e.SetProgramCounter(operand, true)
	return true, nil
}

// JumpZero (JZ) Loads the program counter from the U field *IF* Aa is zero
func JumpZero(e *InstructionEngine) (bool, pkg.Interrupt) {
	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsZero() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// DoublePrecisionJumpZero (DJZ) Loads the program counter from the U field *IF* Aa | Aa+1 is zero
func DoublePrecisionJumpZero(e *InstructionEngine) (bool, pkg.Interrupt) {
	aReg1 := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA()).GetW()
	aReg2 := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA() + 1).GetW()

	if ((aReg1 == pkg.PositiveZero) && (aReg2 == pkg.PositiveZero)) ||
		((aReg1 == pkg.NegativeZero) && (aReg2 == pkg.NegativeOne)) {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNonZero (JNZ) Loads the program counter from the U field *IF* Aa is not zero
func JumpNonZero(e *InstructionEngine) (bool, pkg.Interrupt) {
	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if !aReg.IsZero() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpPositive (JP) Loads the program counter from the U field *IF* Aa is positive
func JumpPositive(e *InstructionEngine) (bool, pkg.Interrupt) {
	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsPositive() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpPositiveAndShift (JPS) Loads the program counter from the U field *IF* Aa is positive.
// Aa is shifted left circularly by one bit in any case
func JumpPositiveAndShift(e *InstructionEngine) (bool, pkg.Interrupt) {
	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsPositive() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	aReg.ShiftLeftCircular(1)
	return true, nil
}

// JumpNegative (JN) Loads the program counter from the U field *IF* Aa is negative
func JumpNegative(e *InstructionEngine) (bool, pkg.Interrupt) {
	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsNegative() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNegativeAndShift (JNS) Loads the program counter from the U field *IF* Aa is negative.
// Aa is shifted left circularly by one bit in any case
func JumpNegativeAndShift(e *InstructionEngine) (bool, pkg.Interrupt) {
	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsNegative() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	aReg.ShiftLeftCircular(1)
	return true, nil
}

// JumpLowBit (JB) Loads the program counter from the U field *IF* the least significant bit of Aa is set
func JumpLowBit(e *InstructionEngine) (bool, pkg.Interrupt) {
	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if uint64(*aReg)&01 != 0 {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNoLowBit (JNB) Loads the program counter from the U field *IF* the least significant bit of Aa is clear
func JumpNoLowBit(e *InstructionEngine) (bool, pkg.Interrupt) {
	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if uint64(*aReg)&01 == 0 {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpGreaterAndDecrement (JGD) Loads the program counter from the U field *IF* the GRS register indicated by
// the value created by the right-most seven bits of the concatenation of the j-field to the a-field is > 0.
// In any case, the indicated register is decremented after the comparison to zero.
func JumpGreaterAndDecrement(e *InstructionEngine) (bool, pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	ix := ((ci.GetJ() << 4) | ci.GetA()) & 0177
	value := e.generalRegisterSet.registers[ix]
	if value.IsPositive() && !value.IsZero() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}

	newValue := pkg.AddSimple(value.GetW(), pkg.NegativeOne)
	e.generalRegisterSet.registers[ix].SetW(newValue)
	return true, nil
}

// JumpModifierGreaterAndIncrement (JMGI) Loads the program counter from the U field *IF* the modifier portion of Xa
// is greater than zero. In any case, the signed increment of Xa is added to the signed modifier of Xa
// (after the comparison to zero).
func JumpModifierGreaterAndIncrement(e *InstructionEngine) (bool, pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	xReg := e.GetExecOrUserXRegister(ci.GetA())
	reg := pkg.Word36(*xReg)
	modifier := reg.GetXH2()
	if pkg.IsPositive(modifier) && !pkg.IsZero(modifier) {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}

	xReg.IncrementModifier()
	return true, nil
}

// JumpOverflow (JO) Loads the program counter from the U field *IF* designator register bit 19 is set.
func JumpOverflow(e *InstructionEngine) (bool, pkg.Interrupt) {
	if e.activityStatePacket.GetDesignatorRegister().IsOverflowSet() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNoOverflow (JNO) Loads the program counter from the U field *IF* designator register bit 19 is clear.
func JumpNoOverflow(e *InstructionEngine) (bool, pkg.Interrupt) {
	if !e.activityStatePacket.GetDesignatorRegister().IsOverflowSet() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpCarry (JC) Loads the program counter from the U field *IF* designator register bit 18 is set.
func JumpCarry(e *InstructionEngine) (bool, pkg.Interrupt) {
	if e.activityStatePacket.GetDesignatorRegister().IsCarrySet() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNoCarry (JNC) Loads the program counter from the U field *IF* designator register bit 18 is clear.
func JumpNoCarry(e *InstructionEngine) (bool, pkg.Interrupt) {
	if !e.activityStatePacket.GetDesignatorRegister().IsCarrySet() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpDivideFault (JDF) Loads the program counter from the U field *IF* designator register bit 23 is set.
func JumpDivideFault(e *InstructionEngine) (bool, pkg.Interrupt) {
	if e.activityStatePacket.GetDesignatorRegister().IsDivideCheckSet() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNoDivideFault (JNDF) Loads the program counter from the U field *IF* designator register bit 23 is clear.
func JumpNoDivideFault(e *InstructionEngine) (bool, pkg.Interrupt) {
	if !e.activityStatePacket.GetDesignatorRegister().IsDivideCheckSet() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpFloatingOverflow (JFO) Loads the program counter from the U field *IF* designator register bit 22 is set.
func JumpFloatingOverflow(e *InstructionEngine) (bool, pkg.Interrupt) {
	if e.activityStatePacket.GetDesignatorRegister().IsCharacteristicOverflowSet() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNoFloatingOverflow (JNFO) Loads the program counter from the U field *IF* designator register bit 22 is clear.
func JumpNoFloatingOverflow(e *InstructionEngine) (bool, pkg.Interrupt) {
	if !e.activityStatePacket.GetDesignatorRegister().IsCharacteristicOverflowSet() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpFloatingUnderflow (JFU) Loads the program counter from the U field *IF* designator register bit 21 is set.
func JumpFloatingUnderflow(e *InstructionEngine) (bool, pkg.Interrupt) {
	if e.activityStatePacket.GetDesignatorRegister().IsCharacteristicUnderflowSet() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNoFloatingUnderflow (JNFU) Loads the program counter from the U field *IF* designator register bit 21 is clear.
func JumpNoFloatingUnderflow(e *InstructionEngine) (bool, pkg.Interrupt) {
	if !e.activityStatePacket.GetDesignatorRegister().IsCharacteristicUnderflowSet() {
		completed, operand, interrupt := e.GetJumpOperand(true)
		if !completed || interrupt != nil {
			return false, interrupt
		}
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}
