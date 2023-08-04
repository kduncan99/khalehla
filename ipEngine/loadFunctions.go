// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

//	TODO DoubleLoadAccumulator (DL)
//	TODO DoubleLoadMagnitudeA (DLM)
//	TODO DoubleLoadNegativeA (DLN)

// LoadAccumulator (LA) loads the content of U under j-field control, and stores it in A(a)
func LoadAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserARegister(ci.GetA()).SetW(operand)
	return true, nil
}

//	TODO LoadAQuarterWord (LAQW)

// LoadIndexRegister (LX) loads the content of U under j-field control, and stores it in X(a)
func LoadIndexRegister(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserXRegister(ci.GetA()).SetW(operand)
	return true, nil
}

// LoadIndexRegisterModifier (LXM)
func LoadIndexRegisterModifier(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserXRegister(ci.GetA()).SetXM(operand)
	return true, nil
}

// LoadIndexRegisterLongModifier (LXLM)
func LoadIndexRegisterLongModifier(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	dr := e.activityStatePacket.GetDesignatorRegister()
	if dr.IsBasicModeEnabled() && dr.GetProcessorPrivilege() > 0 {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
	}

	completed, operand, interrupt := e.GetOperand(true, true, false, false)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserXRegister(ci.GetA()).SetXM24(operand)
	return true, nil
}

// LoadIndexRegisterIncrement (LXI)
func LoadIndexRegisterIncrement(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserXRegister(ci.GetA()).SetXI(operand)
	return true, nil
}

// LoadIndexRegisterShortIncrement (LXSI)
func LoadIndexRegisterShortIncrement(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserXRegister(ci.GetA()).SetXI12(operand)
	return true, nil
}

// LoadMagnitudeAccumulator (LMA)
func LoadMagnitudeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	if pkg.IsNegative(operand) {
		operand ^= pkg.NegativeZero
	}
	ci := e.GetCurrentInstruction()
	e.GetExecOrUserARegister(ci.GetA()).SetW(operand)
	return true, nil
}

// LoadNegativeAccumulator (LNA)
func LoadNegativeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	operand ^= pkg.NegativeZero
	ci := e.GetCurrentInstruction()
	e.GetExecOrUserARegister(ci.GetA()).SetW(operand)
	return true, nil
}

// LoadNegativeMagnitudeAccumulator (LNMA)
func LoadNegativeMagnitudeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	if !pkg.IsNegative(operand) {
		operand ^= pkg.NegativeZero
	}
	ci := e.GetCurrentInstruction()
	e.GetExecOrUserARegister(ci.GetA()).SetW(operand)
	return true, nil
}

// LoadRegister (LR) loads the content of U under j-field control, and stores it in R(a)
func LoadRegister(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserRRegister(ci.GetA()).SetW(operand)
	return true, nil
}

//	TODO LoadRegisterSet (LRS)
//	TODO LoadStringBitLength (LSBL)
//	TODO LoadStringBitOffset (LSBO)
