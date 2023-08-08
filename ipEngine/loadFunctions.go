// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

// DoubleLoadAccumulator (DL) loads the content of U and U+1, storing the values in Aa and Aa+1
func DoubleLoadAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operands, interrupt := e.GetConsecutiveOperands(true, 2, false)
	if !completed || interrupt != nil {
		return
	}

	ci := e.GetCurrentInstruction()
	grsIndex := e.GetExecOrUserARegisterIndex(ci.GetA())
	e.generalRegisterSet.SetRegisterValue(grsIndex, operands[0])
	e.generalRegisterSet.SetRegisterValue(grsIndex+1, operands[1])

	return
}

// DoubleLoadMagnitudeAccumulator (DL) loads the arithmetic magnitude of the content of U and U+1,
// storing the values in Aa and Aa+1
func DoubleLoadMagnitudeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operands, interrupt := e.GetConsecutiveOperands(true, 2, false)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		grsIndex := e.GetExecOrUserARegisterIndex(ci.GetA())
		if operands[0].IsNegative() {
			e.generalRegisterSet.SetRegisterValue(grsIndex, pkg.Word36(pkg.Not(operands[0].GetW())))
			e.generalRegisterSet.SetRegisterValue(grsIndex+1, pkg.Word36(pkg.Not(operands[1].GetW())))
		} else {
			e.generalRegisterSet.SetRegisterValue(grsIndex, operands[0])
			e.generalRegisterSet.SetRegisterValue(grsIndex+1, operands[1])
		}
	}

	return
}

// DoubleLoadNegativeAccumulator (DL) loads the arithmetic negative of the content of U and U+1,
// storing the values in Aa and Aa+1
func DoubleLoadNegativeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operands, interrupt := e.GetConsecutiveOperands(true, 2, false)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		grsIndex := e.GetExecOrUserARegisterIndex(ci.GetA())
		e.generalRegisterSet.SetRegisterValue(grsIndex, pkg.Word36(pkg.Not(operands[0].GetW())))
		e.generalRegisterSet.SetRegisterValue(grsIndex+1, pkg.Word36(pkg.Not(operands[1].GetW())))
	}

	return
}

// LoadAccumulator (LA) loads the content of U under j-field control, and stores it in A(a)
func LoadAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(operand)
	}

	return
}

//	TODO LoadAQuarterWord (LAQW)

// LoadIndexRegister (LX) loads the content of U under j-field control, and stores it in X(a)
func LoadIndexRegister(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetW(operand)
	}

	return
}

// LoadIndexRegisterModifier (LXM)
func LoadIndexRegisterModifier(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetXM(operand)
	}

	return
}

// LoadIndexRegisterLongModifier (LXLM)
func LoadIndexRegisterLongModifier(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = false
	dr := e.activityStatePacket.GetDesignatorRegister()
	if dr.IsBasicModeEnabled() && dr.GetProcessorPrivilege() > 0 {
		interrupt = pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
		return
	}

	completed, operand, interrupt := e.GetOperand(true, true, false, false)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetXM24(operand)
	}

	return
}

// LoadIndexRegisterIncrement (LXI)
func LoadIndexRegisterIncrement(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetXI(operand)
	}

	return
}

// LoadIndexRegisterShortIncrement (LXSI)
func LoadIndexRegisterShortIncrement(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetXI12(operand)
	}

	return
}

// LoadMagnitudeAccumulator (LMA)
func LoadMagnitudeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if completed && interrupt == nil {
		if pkg.IsNegative(operand) {
			operand ^= pkg.NegativeZero
		}
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(operand)
	}

	return
}

// LoadNegativeAccumulator (LNA)
func LoadNegativeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if completed && interrupt == nil {
		operand ^= pkg.NegativeZero
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(operand)
	}

	return
}

// LoadNegativeMagnitudeAccumulator (LNMA)
func LoadNegativeMagnitudeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if completed && interrupt == nil {
		if !pkg.IsNegative(operand) {
			operand ^= pkg.NegativeZero
		}
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(operand)
	}

	return
}

// LoadRegister (LR) loads the content of U under j-field control, and stores it in R(a)
func LoadRegister(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, operand, interrupt := e.GetOperand(true, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserRRegister(ci.GetA()).SetW(operand)
	}

	return
}

//	TODO LoadRegisterSet (LRS)
//	TODO LoadStringBitLength (LSBL)
//	TODO LoadStringBitOffset (LSBO)
