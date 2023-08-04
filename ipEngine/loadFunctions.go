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

//	TODO LoadIndexRegisterModifier (LXM)
//	TODO LoadIndexRegisterLongModifier (LXLM)
//	TODO LoadIndexRegisterIncrement (LXI)
//	TODO LoadIndexRegisterShortIncrement (LXSI)
//	TODO LoadMagnitudeAccumulator (LMA)
//	TODO LoadNegativeAccumulator (LNA)
//	TODO LoadNegativeMagnitudeAccumulator (LNMA)

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
