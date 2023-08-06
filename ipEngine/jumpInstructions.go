// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

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

//	TODO Jump Keys (JK)

// Jump (J) Loads the program counter from the U field - assumes no bank switching

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

//	TODO Halt Keys and Jump (HKJ)
//	TODO Halt Jump (HLTJ)

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
//	TODO Jump Overflow (JO)
//	TODO Jump No Overflow (JNO)
//	TODO Jump Carry (JC)
//	TODO Jump No Carry (JNC)
//	TODO Jump Divide Fault (JDF)
//	TODO Jump No Divide Fault (JNDF)
//	TODO Jump Floating Overflow (JFO)
//	TODO Jump No Floating Overflow (JNFO)
//	TODO Jump Floating Underflow (JFU)
//	TODO Jump No Floating Underflow (JNFU)
