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

	count = count & 0177
	ci := e.GetCurrentInstruction()
	aReg := e.GetExecOrUserARegister(ci.GetA())
	aReg.SetW(pkg.RightShiftCircular(aReg.GetW(), count))
	return true
}

// DoubleShiftCircular (SSC)
func DoubleShiftCircular(e *InstructionEngine) (completed bool) {
	count, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	count = count & 0177
	ci := e.GetCurrentInstruction()
	aReg0 := e.GetExecOrUserARegister(ci.GetA())
	aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
	operand := []uint64{aReg0.GetW(), aReg1.GetW()}
	result := pkg.RightDoubleShiftCircular(operand, count)
	aReg0.SetW(result[0])
	aReg1.SetW(result[1])
	return true
}

// SingleShiftLogical (SSL)
func SingleShiftLogical(e *InstructionEngine) (completed bool) {
	count, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	count = count & 0177
	ci := e.GetCurrentInstruction()
	aReg := e.GetExecOrUserARegister(ci.GetA())
	aReg.SetW(pkg.RightShiftLogical(aReg.GetW(), count))
	return true
}

// DoubleShiftLogical (DSL)
func DoubleShiftLogical(e *InstructionEngine) (completed bool) {
	count, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	count = count & 0177
	ci := e.GetCurrentInstruction()
	aReg0 := e.GetExecOrUserARegister(ci.GetA())
	aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
	operand := []uint64{aReg0.GetW(), aReg1.GetW()}
	result := pkg.RightDoubleShiftLogical(operand, count)
	aReg0.SetW(result[0])
	aReg1.SetW(result[1])
	return true
}

// SingleShiftAlgebraic (SSA)
func SingleShiftAlgebraic(e *InstructionEngine) (completed bool) {
	count, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	count = (count & 0177) % 36
	ci := e.GetCurrentInstruction()
	aReg := e.GetExecOrUserARegister(ci.GetA())
	aReg.SetW(pkg.RightShiftAlgebraic(aReg.GetW(), count))
	return true
}

// DoubleShiftAlgebraic (DSA)
func DoubleShiftAlgebraic(e *InstructionEngine) (completed bool) {
	count, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	count = count & 0177
	ci := e.GetCurrentInstruction()
	aReg0 := e.GetExecOrUserARegister(ci.GetA())
	aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
	operand := []uint64{aReg0.GetW(), aReg1.GetW()}
	result := pkg.RightDoubleShiftAlgebraic(operand, count)
	aReg0.SetW(result[0])
	aReg1.SetW(result[1])
	return true
}

func bitsMatch(value uint64) bool {
	return ((value >> 35) & 01) == ((value >> 34) & 01)
}

// LoadShiftAndCount (LSC) performs a circular shift left of U until bit 0 != bit 1, storing the number of shifts
// performed into Aa+1. If bit 0 starts out != bit 1, no shift is performed.
// In either case, the value of U after shifting (possibly 0 times) is stored in Aa.
// If U contains +/- zero then the value is stored in Aa, and the shift count stored in Aa+1 will be 35.
// This can be used for scaling / normalizing a floating point number.
func LoadShiftAndCount(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(false, true, false, false, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg0 := e.GetExecOrUserARegister(ci.GetA())
		aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
		value := result.operand
		count := 0

		if pkg.IsZero(value) {
			count = 35
		} else {
			for bitsMatch(value) {
				partial := value & 01
				value >>= 1
				value |= partial << 35
				count++
			}
		}

		aReg0.SetW(value)
		aReg1.SetW(uint64(count))
	}

	return result.complete
}

// DoubleLoadShiftAndCount (DLSC) performs a circular shift left of U/U+1 until bit 0 != bit 1, storing the number of
// shifts performed into Aa+2. If bit 0 starts out != bit 1, no shift is performed.
// In either case, the value of U after shifting (possibly 0 times) is stored in Aa/Aa+1.
// If U/U+1 contains +/- zero then the value is stored in Aa/Aa+1, and the shift count stored in Aa+2 will be 71.
// This can be used for scaling / normalizing a floating point number.
func DoubleLoadShiftAndCount(e *InstructionEngine) (completed bool) {
	result := e.GetConsecutiveOperands(true, 2, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg0 := e.GetExecOrUserARegister(ci.GetA())
		aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
		aReg2 := e.GetExecOrUserARegister(ci.GetA() + 2)
		value := []uint64{result.source[0].GetW(), result.source[1].GetW()}
		count := 0

		if pkg.IsDoubleZero(value) {
			count = 71
		} else {
			for bitsMatch(value[0]) {
				partial0 := value[0] & 01
				partial1 := value[1] & 01
				value[0] >>= 1
				value[1] >>= 1
				value[0] |= partial0 << 35
				value[1] |= partial1 << 35
				count++
			}
		}

		aReg0.SetW(value[0])
		aReg1.SetW(value[1])
		aReg2.SetW(uint64(count))
	}

	return result.complete
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

// LeftDoubleShiftCircular (LDSC)
func LeftDoubleShiftCircular(e *InstructionEngine) (completed bool) {
	count, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	count = count & 0177
	ci := e.GetCurrentInstruction()
	aReg0 := e.GetExecOrUserARegister(ci.GetA())
	aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
	operand := []uint64{aReg0.GetW(), aReg1.GetW()}
	result := pkg.LeftDoubleShiftCircular(operand, count)
	aReg0.SetW(result[0])
	aReg1.SetW(result[1])
	return true
}

// LeftSingleShiftLogical (LSSL)
func LeftSingleShiftLogical(e *InstructionEngine) (completed bool) {
	count, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	count = (count & 0177) % 36
	ci := e.GetCurrentInstruction()
	aReg := e.GetExecOrUserARegister(ci.GetA())
	aReg.SetW(pkg.LeftShiftLogical(aReg.GetW(), count))
	return true
}

// LeftDoubleShiftLogical (LDSL)
func LeftDoubleShiftLogical(e *InstructionEngine) (completed bool) {
	count, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	count = count & 0177
	ci := e.GetCurrentInstruction()
	aReg0 := e.GetExecOrUserARegister(ci.GetA())
	aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
	operand := []uint64{aReg0.GetW(), aReg1.GetW()}
	result := pkg.LeftDoubleShiftLogical(operand, count)
	aReg0.SetW(result[0])
	aReg1.SetW(result[1])
	return true
}
