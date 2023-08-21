// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

// SingleShiftCircular (SSC) shifts the value in Aa to the right, by the number of bits indicated in U.
// bits shifted out of bit 35 are shifted into bit 0.
func SingleShiftCircular(e *InstructionEngine) (completed bool) {
	op, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	op = op % 36
	ci := e.GetCurrentInstruction()
	aReg := e.GetExecOrUserARegister(ci.GetA())
	value := aReg.GetW()
	if op > 18 {
		remnant := value & 0_777777
		value >>= 18
		value |= remnant << 18
		op -= 18
	}

	if op > 0 {
		mask := (2 ^ op) - 1
		remnant := value & mask
		value >>= op
		value |= remnant << (36 - op)
	}

	aReg.SetW(value)
	return true
}

// DoubleShiftCircular (SSC)
func DoubleShiftCircular(e *InstructionEngine) (completed bool) {
	//	TODO
	return true
}

// SingleShiftLogical (SSL)
func SingleShiftLogical(e *InstructionEngine) (completed bool) {
	//	TODO
	return true
}

// DoubleShiftLogical (DSL)
func DoubleShiftLogical(e *InstructionEngine) (completed bool) {
	//	TODO
	return true
}

// SingleShiftAlgebraic (SSA)
func SingleShiftAlgebraic(e *InstructionEngine) (completed bool) {
	//	TODO
	return true
}

// DoubleShiftAlgebraic (DSA)
func DoubleShiftAlgebraic(e *InstructionEngine) (completed bool) {
	//	TODO
	return true
}

// LoadShiftAndCount (LSC)
func LoadShiftAndCount(e *InstructionEngine) (completed bool) {
	//	TODO
	return true
}

// DoubleLoadShiftAndCount (DLSC)
func DoubleLoadShiftAndCount(e *InstructionEngine) (completed bool) {
	//	TODO
	return true
}

// LeftSingleShiftCircular (LSSC) shifts the value in Aa to the left, by the number of bits indicated in U.
// bits shifted out of bit 0 are shifted into bit 35.
func LeftSingleShiftCircular(e *InstructionEngine) (completed bool) {
	op, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	op = op % 36
	ci := e.GetCurrentInstruction()
	aReg := e.GetExecOrUserARegister(ci.GetA())
	value := aReg.GetW()
	if op > 18 {
		value <<= 18
		value |= value >> 36
		value &= pkg.NegativeZero
	}

	if op > 0 {
		value <<= op
		value |= value >> 36
	}

	aReg.SetW(value)
	return true
}

// LeftDoubleShiftCircular (LSSC)
func LeftDoubleShiftCircular(e *InstructionEngine) (completed bool) {
	//	TODO
	return true
}

// LeftSingleShiftLogical (LSSL)
func LeftSingleShiftLogical(e *InstructionEngine) (completed bool) {
	//	TODO
	return true
}

// LeftDoubleShiftLogical (LDSL)
func LeftDoubleShiftLogical(e *InstructionEngine) (completed bool) {
	//	TODO
	return true
}
