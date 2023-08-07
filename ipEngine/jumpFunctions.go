// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

//	TODO Store Location and Jump (SLJ)
//	TODO Load Modifier and Jump (LMJ)

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

//	TODO Jump Zero (JZ)
//	TODO Double-Precision Jump Zero (DJZ)
//	TODO Jump Nonzero (JNZ)
//	TODO Jump Positive (JP)
//	TODO Jump Positive and Shift (JPS)
//	TODO Jump Negative (JN)
//	TODO Jump Negative and Shift (JNS)
//	TODO Jump Low Bit (JB)
//	TODO Jump No Low Bit (JNB)
//	TODO Jump Greater and Decrement (JGD)
//	TODO Jump Modifier Greater and Increment (JMGI)

// JumpOverflow (JO) Loads the program counter from the U field *IF* designator register bit 19 is set.
func JumpOverflow(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	if e.activityStatePacket.GetDesignatorRegister().IsOverflowSet() {
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNoOverflow (JNO) Loads the program counter from the U field *IF* designator register bit 19 is clear.
func JumpNoOverflow(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	if !e.activityStatePacket.GetDesignatorRegister().IsOverflowSet() {
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpCarry (JC) Loads the program counter from the U field *IF* designator register bit 18 is set.
func JumpCarry(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	if e.activityStatePacket.GetDesignatorRegister().IsCarrySet() {
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNoCarry (JNC) Loads the program counter from the U field *IF* designator register bit 18 is clear.
func JumpNoCarry(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	if !e.activityStatePacket.GetDesignatorRegister().IsCarrySet() {
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpDivideFault (JDF) Loads the program counter from the U field *IF* designator register bit 23 is set.
func JumpDivideFault(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	if e.activityStatePacket.GetDesignatorRegister().IsDivideCheckSet() {
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNoDivideFault (JNDF) Loads the program counter from the U field *IF* designator register bit 23 is clear.
func JumpNoDivideFault(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	if !e.activityStatePacket.GetDesignatorRegister().IsDivideCheckSet() {
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpFloatingOverflow (JFO) Loads the program counter from the U field *IF* designator register bit 22 is set.
func JumpFloatingOverflow(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	if e.activityStatePacket.GetDesignatorRegister().IsCharacteristicOverflowSet() {
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNoFloatingOverflow (JNFO) Loads the program counter from the U field *IF* designator register bit 22 is clear.
func JumpNoFloatingOverflow(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	if !e.activityStatePacket.GetDesignatorRegister().IsCharacteristicOverflowSet() {
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpFloatingUnderflow (JFU) Loads the program counter from the U field *IF* designator register bit 21 is set.
func JumpFloatingUnderflow(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	if e.activityStatePacket.GetDesignatorRegister().IsCharacteristicUnderflowSet() {
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}

// JumpNoFloatingUnderflow (JNFU) Loads the program counter from the U field *IF* designator register bit 21 is clear.
func JumpNoFloatingUnderflow(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetJumpOperand(true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	if !e.activityStatePacket.GetDesignatorRegister().IsCharacteristicUnderflowSet() {
		e.SetProgramCounter(operand, true)
	}
	return true, nil
}
