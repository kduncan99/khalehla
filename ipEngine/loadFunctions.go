// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

// DoubleLoadAccumulator (DL) loads the content of U and U+1, storing the values in Aa and Aa+1
func DoubleLoadAccumulator(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operands, interrupt := e.GetConsecutiveOperands(true, 2, false)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	grsIndex := e.GetExecOrUserARegisterIndex(ci.GetA())
	e.generalRegisterSet.SetRegisterValue(grsIndex, operands[0])
	e.generalRegisterSet.SetRegisterValue(grsIndex+1, operands[1])

	return true, nil
}

// DoubleLoadMagnitudeAccumulator (DL) loads the arithmetic magnitude of the content of U and U+1,
// storing the values in Aa and Aa+1
func DoubleLoadMagnitudeAccumulator(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operands, interrupt := e.GetConsecutiveOperands(true, 2, false)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	grsIndex := e.GetExecOrUserARegisterIndex(ci.GetA())
	if operands[0].IsNegative() {
		e.generalRegisterSet.SetRegisterValue(grsIndex, pkg.Word36(pkg.Not(operands[0].GetW())))
		e.generalRegisterSet.SetRegisterValue(grsIndex+1, pkg.Word36(pkg.Not(operands[1].GetW())))
	} else {
		e.generalRegisterSet.SetRegisterValue(grsIndex, operands[0])
		e.generalRegisterSet.SetRegisterValue(grsIndex+1, operands[1])
	}

	return true, nil
}

// DoubleLoadNegativeAccumulator (DL) loads the arithmetic negative of the content of U and U+1,
// storing the values in Aa and Aa+1
func DoubleLoadNegativeAccumulator(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operands, interrupt := e.GetConsecutiveOperands(true, 2, false)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	grsIndex := e.GetExecOrUserARegisterIndex(ci.GetA())
	e.generalRegisterSet.SetRegisterValue(grsIndex, pkg.Word36(pkg.Not(operands[0].GetW())))
	e.generalRegisterSet.SetRegisterValue(grsIndex+1, pkg.Word36(pkg.Not(operands[1].GetW())))

	return true, nil
}

// LoadAccumulator (LA) loads the content of U under j-field control, and stores it in A(a)
func LoadAccumulator(e *InstructionEngine) (bool, pkg.Interrupt) {
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
func LoadIndexRegister(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserXRegister(ci.GetA()).SetW(operand)
	return true, nil
}

// LoadIndexRegisterModifier (LXM)
func LoadIndexRegisterModifier(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserXRegister(ci.GetA()).SetXM(operand)
	return true, nil
}

// LoadIndexRegisterLongModifier (LXLM)
func LoadIndexRegisterLongModifier(e *InstructionEngine) (bool, pkg.Interrupt) {
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
func LoadIndexRegisterIncrement(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserXRegister(ci.GetA()).SetXI(operand)
	return true, nil
}

// LoadIndexRegisterShortIncrement (LXSI)
func LoadIndexRegisterShortIncrement(e *InstructionEngine) (bool, pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if !completed || interrupt != nil {
		return false, interrupt
	}

	ci := e.GetCurrentInstruction()
	e.GetExecOrUserXRegister(ci.GetA()).SetXI12(operand)
	return true, nil
}

// LoadMagnitudeAccumulator (LMA)
func LoadMagnitudeAccumulator(e *InstructionEngine) (bool, pkg.Interrupt) {
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
func LoadNegativeAccumulator(e *InstructionEngine) (bool, pkg.Interrupt) {
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
func LoadNegativeMagnitudeAccumulator(e *InstructionEngine) (bool, pkg.Interrupt) {
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
func LoadRegister(e *InstructionEngine) (bool, pkg.Interrupt) {
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
