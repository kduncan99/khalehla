// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

// LogicalOr performs a bit-wise OR on the content of (U) and Aa, storing the result into Aa+1
func LogicalOr(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(false, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		value := result.operand | e.GetExecOrUserARegister(ci.GetA()).GetW()
		e.GetExecOrUserARegister(ci.GetA() + 1).SetW(value)
	}

	return true
}

// LogicalExclusiveOr performs a bit-wise XOR on the content of (U) and Aa, storing the result into Aa+1
func LogicalExclusiveOr(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(false, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		value := result.operand ^ e.GetExecOrUserARegister(ci.GetA()).GetW()
		e.GetExecOrUserARegister(ci.GetA() + 1).SetW(value)
	}

	return true
}

// LogicalAnd performs a bit-wise AND on the content of (U) and Aa, storing the result into Aa+1
func LogicalAnd(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(false, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		value := result.operand & e.GetExecOrUserARegister(ci.GetA()).GetW()
		e.GetExecOrUserARegister(ci.GetA() + 1).SetW(value)
	}

	return true
}

// MaskedLoadUpper produces the value resulting from a bit-wise OR of (U) AND R2, and Aa AND NOT R2,
// storing that value into Aa+1.
func MaskedLoadUpper(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(false, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserXRegister(ci.GetA()).GetW()
		mask := e.GetExecOrUserRRegister(pkg.R2).GetW()
		notMask := pkg.Not(mask)
		value := (result.operand & mask) | (aValue & notMask)
		e.GetExecOrUserARegister(ci.GetA() + 1).SetW(value)
	}

	return true
}
