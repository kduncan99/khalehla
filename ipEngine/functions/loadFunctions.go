// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package functions

import "khalehla/ipEngine"

//	TODO DoubleLoadAccumulator (DL)
//	TODO DoubleLoadMagnitudeA (DLM)
//	TODO DoubleLoadNegativeA (DLN)

// LoadAccumulator (LA) loads the content of U under j-field control, and stores it in A(a)
func LoadAccumulator(e *ipEngine.InstructionEngine) (completed bool, interrupt ipEngine.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserARegister(uint(ci.GetA())).SetW(operand)
	return true, nil
}

//	TODO LoadAQuarterWord (LAQW)

// LoadIndexRegister (LX) loads the content of U under j-field control, and stores it in X(a)
func LoadIndexRegister(e *ipEngine.InstructionEngine) (completed bool, interrupt ipEngine.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserXRegister(uint(ci.GetA())).SetW(operand)
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
func LoadRegister(e *ipEngine.InstructionEngine) (completed bool, interrupt ipEngine.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserRRegister(uint(ci.GetA())).SetW(operand)
	return true, nil
}

//	TODO LoadRegisterSet (LRS)
//	TODO LoadStringBitLength (LSBL)
//	TODO LoadStringBitOffset (LSBO)
