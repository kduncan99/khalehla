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
		flipDesignatorRegisterBit31(e)
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
			flipDesignatorRegisterBit31(e)
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
			flipDesignatorRegisterBit31(e)
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
			flipDesignatorRegisterBit31(e)
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
			flipDesignatorRegisterBit31(e)
		}
	}

	return comp
}

// Jumps conditional upon some value -----------------------------------------------------------------------------------

// JumpZero (JZ) Loads the program counter from the U field *IF* Aa is zero
func JumpZero(e *InstructionEngine) (completed bool) {
	completed = true
	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())

	if aReg.IsZero() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// DoubleJumpZero (DJZ) Loads the program counter from the U field *IF* Aa | Aa+1 is zero
func DoubleJumpZero(e *InstructionEngine) (completed bool) {
	completed = true

	ci := e.GetCurrentInstruction()
	ax := ci.GetA()
	aReg1 := e.GetExecOrUserARegister(ax).GetW()
	aReg2 := e.GetExecOrUserARegister(ax + 1).GetW()

	if ((aReg1 == pkg.PositiveZero) && (aReg2 == pkg.PositiveZero)) ||
		((aReg1 == pkg.NegativeZero) && (aReg2 == pkg.NegativeZero)) {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpNonZero (JNZ) Loads the program counter from the U field *IF* Aa is not zero
func JumpNonZero(e *InstructionEngine) (completed bool) {
	completed = true

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if !aReg.IsZero() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpPositive (JP) Loads the program counter from the U field *IF* Aa is positive
func JumpPositive(e *InstructionEngine) (completed bool) {
	completed = true

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsPositive() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpPositiveAndShift (JPS) Loads the program counter from the U field *IF* Aa is positive.
// Aa is shifted left circularly by one bit in any case
func JumpPositiveAndShift(e *InstructionEngine) (completed bool) {
	completed = true

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsPositive() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	if completed {
		aReg.ShiftLeftCircular(1)
	}

	return
}

// JumpNegative (JN) Loads the program counter from the U field *IF* Aa is negative
func JumpNegative(e *InstructionEngine) (completed bool) {
	completed = true

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsNegative() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpNegativeAndShift (JNS) Loads the program counter from the U field *IF* Aa is negative.
// Aa is shifted left circularly by one bit in any case
func JumpNegativeAndShift(e *InstructionEngine) (completed bool) {
	completed = true

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if aReg.IsNegative() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	if completed {
		aReg.ShiftLeftCircular(1)
	}

	return
}

// JumpLowBit (JB) Loads the program counter from the U field *IF* the least significant bit of Aa is set
func JumpLowBit(e *InstructionEngine) (completed bool) {
	completed = true

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if uint64(*aReg)&01 != 0 {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpNoLowBit (JNB) Loads the program counter from the U field *IF* the least significant bit of Aa is clear
func JumpNoLowBit(e *InstructionEngine) (completed bool) {
	completed = true

	aReg := e.GetExecOrUserARegister(e.GetCurrentInstruction().GetA())
	if uint64(*aReg)&01 == 0 {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpGreaterAndDecrement (JGD) Loads the program counter from the U field *IF* the GRS register indicated by
// the value created by the right-most seven bits of the concatenation of the j-field to the a-field is > 0.
// In any case, the indicated register is decremented after the comparison to zero.
func JumpGreaterAndDecrement(e *InstructionEngine) (completed bool) {
	completed = true

	ci := e.GetCurrentInstruction()
	ix := ((ci.GetJ() << 4) | ci.GetA()) & 0177
	value := e.generalRegisterSet.registers[ix]
	if value.IsPositive() && !value.IsZero() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	if completed {
		newValue := pkg.AddSimple(value.GetW(), pkg.NegativeOne)
		e.generalRegisterSet.registers[ix].SetW(newValue)
	}

	return
}

// JumpModifierGreaterAndIncrement (JMGI) Loads the program counter from the U field *IF* the modifier portion of Xa
// is greater than zero. In any case, the signed increment of Xa is added to the signed modifier of Xa
// (after the comparison to zero).
func JumpModifierGreaterAndIncrement(e *InstructionEngine) (completed bool) {
	completed = true

	ci := e.GetCurrentInstruction()
	xReg := e.GetExecOrUserXRegister(ci.GetA())
	reg := pkg.Word36(*xReg)
	modifier := reg.GetXH2()
	if pkg.IsPositive(modifier) && !pkg.IsZero(modifier) {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	if completed {
		xReg.IncrementModifier()
	}

	return
}

// Jumps conditional upon a designator register bit --------------------------------------------------------------------

// JumpOverflow (JO) Loads the program counter from the U field *IF* designator register bit 19 is set.
func JumpOverflow(e *InstructionEngine) (completed bool) {
	completed = true

	if e.activityStatePacket.GetDesignatorRegister().IsOverflowSet() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpNoOverflow (JNO) Loads the program counter from the U field *IF* designator register bit 19 is clear.
func JumpNoOverflow(e *InstructionEngine) (completed bool) {
	completed = true

	if !e.activityStatePacket.GetDesignatorRegister().IsOverflowSet() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpCarry (JC) Loads the program counter from the U field *IF* designator register bit 18 is set.
func JumpCarry(e *InstructionEngine) (completed bool) {
	completed = true

	if e.activityStatePacket.GetDesignatorRegister().IsCarrySet() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpNoCarry (JNC) Loads the program counter from the U field *IF* designator register bit 18 is clear.
func JumpNoCarry(e *InstructionEngine) (completed bool) {
	completed = true

	if !e.activityStatePacket.GetDesignatorRegister().IsCarrySet() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpDivideFault (JDF) Loads the program counter from the U field *IF* designator register bit 23 is set.
func JumpDivideFault(e *InstructionEngine) (completed bool) {
	completed = true

	if e.activityStatePacket.GetDesignatorRegister().IsDivideCheckSet() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpNoDivideFault (JNDF) Loads the program counter from the U field *IF* designator register bit 23 is clear.
func JumpNoDivideFault(e *InstructionEngine) (completed bool) {
	completed = true

	if !e.activityStatePacket.GetDesignatorRegister().IsDivideCheckSet() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpFloatingOverflow (JFO) Loads the program counter from the U field *IF* designator register bit 22 is set.
func JumpFloatingOverflow(e *InstructionEngine) (completed bool) {
	completed = true

	if e.activityStatePacket.GetDesignatorRegister().IsCharacteristicOverflowSet() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpNoFloatingOverflow (JNFO) Loads the program counter from the U field *IF* designator register bit 22 is clear.
func JumpNoFloatingOverflow(e *InstructionEngine) (completed bool) {
	completed = true

	if !e.activityStatePacket.GetDesignatorRegister().IsCharacteristicOverflowSet() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpFloatingUnderflow (JFU) Loads the program counter from the U field *IF* designator register bit 21 is set.
func JumpFloatingUnderflow(e *InstructionEngine) (completed bool) {
	completed = true

	if e.activityStatePacket.GetDesignatorRegister().IsCharacteristicUnderflowSet() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}

// JumpNoFloatingUnderflow (JNFU) Loads the program counter from the U field *IF* designator register bit 21 is clear.
func JumpNoFloatingUnderflow(e *InstructionEngine) (completed bool) {
	completed = true

	if !e.activityStatePacket.GetDesignatorRegister().IsCharacteristicUnderflowSet() {
		operand, flip31, comp, i := e.GetJumpOperand()

		if i != nil {
			e.PostInterrupt(i)
		} else if comp {
			e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
			e.SetProgramCounter(operand, true)
			if flip31 {
				flipDesignatorRegisterBit31(e)
			}
		}

		completed = comp
	}

	return
}
