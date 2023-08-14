// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"
	"khalehla/pkg"
)

// DoubleLoadAccumulator (DL) loads the content of U and U+1, storing the values in Aa and Aa+1
func DoubleLoadAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	result := e.GetConsecutiveOperands(true, 2, false)
	if result.complete && result.interrupt == nil {
		ci := e.GetCurrentInstruction()
		grsIndex := e.GetExecOrUserARegisterIndex(ci.GetA())
		e.generalRegisterSet.SetRegisterValue(grsIndex, result.source[0])
		e.generalRegisterSet.SetRegisterValue(grsIndex+1, result.source[1])
	}

	return result.complete, result.interrupt
}

// DoubleLoadMagnitudeAccumulator (DL) loads the arithmetic magnitude of the content of U and U+1,
// storing the values in Aa and Aa+1
func DoubleLoadMagnitudeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	result := e.GetConsecutiveOperands(true, 2, false)
	if result.complete && result.interrupt == nil {
		ci := e.GetCurrentInstruction()
		grsIndex := e.GetExecOrUserARegisterIndex(ci.GetA())
		if result.source[0].IsNegative() {
			e.generalRegisterSet.SetRegisterValue(grsIndex, pkg.Word36(pkg.Not(result.source[0].GetW())))
			e.generalRegisterSet.SetRegisterValue(grsIndex+1, pkg.Word36(pkg.Not(result.source[1].GetW())))
		} else {
			e.generalRegisterSet.SetRegisterValue(grsIndex, result.source[0])
			e.generalRegisterSet.SetRegisterValue(grsIndex+1, result.source[1])
		}
	}

	return result.complete, result.interrupt
}

// DoubleLoadNegativeAccumulator (DL) loads the arithmetic negative of the content of U and U+1,
// storing the values in Aa and Aa+1
func DoubleLoadNegativeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	result := e.GetConsecutiveOperands(true, 2, false)
	if result.complete && result.interrupt == nil {
		ci := e.GetCurrentInstruction()
		grsIndex := e.GetExecOrUserARegisterIndex(ci.GetA())
		e.generalRegisterSet.SetRegisterValue(grsIndex, pkg.Word36(pkg.Not(result.source[0].GetW())))
		e.generalRegisterSet.SetRegisterValue(grsIndex+1, pkg.Word36(pkg.Not(result.source[1].GetW())))
	}

	return result.complete, result.interrupt
}

// LoadAccumulator (LA) loads the content of U under j-field control, and stores it in A(a)
func LoadAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	result := e.GetOperand(true, true, true, true, false)
	if result.complete && result.interrupt == nil {
		fmt.Printf("FOO") // TODO remove
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(result.operand)
	}

	return result.complete, result.interrupt
}

//	TODO LoadAQuarterWord (LAQW)

// LoadIndexRegister (LX) loads the content of U under j-field control, and stores it in X(a)
func LoadIndexRegister(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	result := e.GetOperand(true, true, true, true, false)
	if result.complete && result.interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetW(result.operand)
	}

	return result.complete, result.interrupt
}

// LoadIndexRegisterModifier (LXM)
func LoadIndexRegisterModifier(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	result := e.GetOperand(true, true, true, true, false)
	if result.complete && result.interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetXM(result.operand)
	}

	return result.complete, result.interrupt
}

// LoadIndexRegisterLongModifier (LXLM)
func LoadIndexRegisterLongModifier(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = false
	dr := e.activityStatePacket.GetDesignatorRegister()
	if dr.IsBasicModeEnabled() && dr.GetProcessorPrivilege() > 0 {
		interrupt = pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
		return
	}

	result := e.GetOperand(true, true, false, false, false)
	if result.complete && result.interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetXM24(result.operand)
	}

	return result.complete, result.interrupt
}

// LoadIndexRegisterIncrement (LXI)
func LoadIndexRegisterIncrement(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	result := e.GetOperand(true, true, true, true, false)
	if result.complete && result.interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetXI(result.operand)
	}

	return result.complete, result.interrupt
}

// LoadIndexRegisterShortIncrement (LXSI)
func LoadIndexRegisterShortIncrement(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	result := e.GetOperand(true, true, true, true, false)
	if result.complete && result.interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetXI12(result.operand)
	}

	return result.complete, result.interrupt
}

// LoadMagnitudeAccumulator (LMA)
func LoadMagnitudeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	result := e.GetOperand(true, true, true, true, false)
	if result.complete && result.interrupt == nil {
		if pkg.IsNegative(result.operand) {
			result.operand ^= pkg.NegativeZero
		}
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(result.operand)
	}

	return result.complete, result.interrupt
}

// LoadNegativeAccumulator (LNA)
func LoadNegativeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	result := e.GetOperand(true, true, true, true, false)
	if result.complete && result.interrupt == nil {
		result.operand ^= pkg.NegativeZero
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(result.operand)
	}

	return result.complete, result.interrupt
}

// LoadNegativeMagnitudeAccumulator (LNMA)
func LoadNegativeMagnitudeAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	result := e.GetOperand(true, true, true, true, false)
	if result.complete && result.interrupt == nil {
		if !pkg.IsNegative(result.operand) {
			result.operand ^= pkg.NegativeZero
		}
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(result.operand)
	}

	return result.complete, result.interrupt
}

// LoadRegister (LR) loads the content of U under j-field control, and stores it in R(a)
func LoadRegister(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	result := e.GetOperand(true, true, true, true, false)
	if result.complete && result.interrupt == nil {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserRRegister(ci.GetA()).SetW(result.operand)
	}

	return result.complete, result.interrupt
}

//	TODO LoadRegisterSet (LRS)
//	TODO LoadStringBitLength (LSBL)
//	TODO LoadStringBitOffset (LSBO)
