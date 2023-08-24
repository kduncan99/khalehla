// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

// SingleShiftCircular (SSC) shifts the value in Aa to the right, by the number of bits indicated in U.
// bits shifted out of bit 35 are shifted into bit 0.
func SingleShiftCircular(e *InstructionEngine) (completed bool) {
	count, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	count = (count & 0177) % 36
	ci := e.GetCurrentInstruction()
	aReg := e.GetExecOrUserARegister(ci.GetA())
	aReg.SetW(pkg.RightShiftCircular(aReg.GetW(), count))
	return true
}

// DoubleShiftCircular (SSC)
func DoubleShiftCircular(e *InstructionEngine) (completed bool) {
	op, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	op = (op & 0177) % 72
	ci := e.GetCurrentInstruction()
	aReg0 := e.GetExecOrUserARegister(ci.GetA())
	aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
	value0 := aReg0.GetW()
	value1 := aReg1.GetW()

	if op > 36 {
		v := value0
		value0 = value1
		value1 = v
		op -= 72
	}

	if op > 18 {
		remnant := value & 0_777777
		value >>= 18
		value |= remnant << 18
		op -= 18
	}

	if op > 0 {
		mask := uint64(2<<op) - 1
		remnant := value & mask
		value >>= op
		value |= remnant << (36 - op)
	}

	aReg.SetW(value)
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
	count, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	count = (count & 0177) % 36
	ci := e.GetCurrentInstruction()
	aReg := e.GetExecOrUserARegister(ci.GetA())
	aReg.SetW(pkg.LeftShiftCircular(aReg.GetW(), count))
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
