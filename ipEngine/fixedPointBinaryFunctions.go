// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"
	"khalehla/pkg"
)

// The following instructions update DB18 (carry) and DB19 (overflow) in the following conditions:
//  Input Signs   Output Sign   DB18 DB19
//      +/+             +         0    0
//      +/+             -         0    1
//      +/-             +         1    0
//      +/-             -         0    0
//      -/-             +         1    1
//      -/-             -         1    0
// AA, ANA, AMA, ANMA, AU, ANU, AX, ANX, DA, DAN, ADD1, SUB1

func updateDesignatorRegister(e *InstructionEngine, addend1Positive bool, addend2Positive bool, sumPositive bool) {
	bothNeg := !addend1Positive && !addend2Positive
	addendsAgree := addend1Positive == addend2Positive
	db18 := bothNeg || (!addendsAgree && sumPositive)
	db19 := addendsAgree && (addend1Positive != sumPositive)
	dr := e.GetDesignatorRegister()
	dr.SetCarry(db18)
	dr.SetOverflow(db19)

	if dr.IsOperationTrapEnabled() && db19 {
		i := pkg.NewOperationTrapInterrupt(pkg.OperationTrapFixedPointBinaryOverflow)
		e.PostInterrupt(i)
	}
}

// AddAccumulator (AA) adds (U) to Aa
func AddAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg := e.GetExecOrUserARegister(ci.GetA())
		addend1 := aReg.GetW()
		addend2 := result.operand
		sum := pkg.AddSimple(addend1, addend2)
		aReg.SetW(sum)
		updateDesignatorRegister(e, pkg.IsPositive(addend1), pkg.IsPositive(addend2), pkg.IsPositive(sum))
	}

	return result.complete
}

// AddNegativeAccumulator (ANA) adds -(U) to Aa
func AddNegativeAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg := e.GetExecOrUserARegister(ci.GetA())
		addend1 := aReg.GetW()
		addend2 := pkg.Negate(result.operand)
		sum := pkg.AddSimple(addend1, addend2)
		aReg.SetW(sum)
		updateDesignatorRegister(e, pkg.IsPositive(addend1), pkg.IsPositive(addend2), pkg.IsPositive(sum))
	}

	return result.complete
}

// AddMagnitudeAccumulator (ANA) adds |(U)| to Aa
func AddMagnitudeAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg := e.GetExecOrUserARegister(ci.GetA())
		addend1 := aReg.GetW()
		addend2 := result.operand
		if pkg.IsNegative(addend2) {
			addend2 = pkg.Negate(result.operand)
		}
		sum := pkg.AddSimple(addend1, addend2)
		aReg.SetW(sum)
		updateDesignatorRegister(e, pkg.IsPositive(addend1), pkg.IsPositive(addend2), pkg.IsPositive(sum))
	}

	return result.complete
}

// AddNegativeMagnitudeAccumulator (ANMA) adds -|(U)| to Aa
func AddNegativeMagnitudeAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg := e.GetExecOrUserARegister(ci.GetA())
		addend1 := aReg.GetW()
		addend2 := result.operand
		if !pkg.IsNegative(addend2) {
			addend2 = pkg.Negate(result.operand)
		}
		sum := pkg.AddSimple(addend1, addend2)
		aReg.SetW(sum)
		updateDesignatorRegister(e, pkg.IsPositive(addend1), pkg.IsPositive(addend2), pkg.IsPositive(sum))
	}

	return result.complete
}

// AddUpperAccumulator (AU) adds (U) and Aa, storing the sum in Aa+1
func AddUpperAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg0 := e.GetExecOrUserARegister(ci.GetA())
		aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
		addend1 := aReg0.GetW()
		addend2 := result.operand
		sum := pkg.AddSimple(addend1, addend2)
		aReg1.SetW(sum)
		updateDesignatorRegister(e, pkg.IsPositive(addend1), pkg.IsPositive(addend2), pkg.IsPositive(sum))
	}

	return result.complete
}

// AddNegativeUpperAccumulator (ANA) adds -(U) and Aa, storing the sum in Aa+1
func AddNegativeUpperAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg0 := e.GetExecOrUserARegister(ci.GetA())
		aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
		addend1 := aReg0.GetW()
		addend2 := pkg.Negate(result.operand)
		sum := pkg.AddSimple(addend1, addend2)
		aReg1.SetW(sum)
		updateDesignatorRegister(e, pkg.IsPositive(addend1), pkg.IsPositive(addend2), pkg.IsPositive(sum))
	}

	return result.complete
}

// AddIndexRegister (AX) adds (U) to Xa
func AddIndexRegister(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		xReg := e.GetExecOrUserXRegister(ci.GetA())
		addend1 := xReg.GetW()
		addend2 := result.operand
		sum := pkg.AddSimple(addend1, addend2)
		xReg.SetW(sum)
		updateDesignatorRegister(e, pkg.IsPositive(addend1), pkg.IsPositive(addend2), pkg.IsPositive(sum))
	}

	return result.complete
}

// AddNegativeIndexRegister (ANX) adds -(U) to Xa
func AddNegativeIndexRegister(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		xReg := e.GetExecOrUserXRegister(ci.GetA())
		addend1 := xReg.GetW()
		addend2 := pkg.Negate(result.operand)
		sum := pkg.AddSimple(addend1, addend2)
		xReg.SetW(sum)
		updateDesignatorRegister(e, pkg.IsPositive(addend1), pkg.IsPositive(addend2), pkg.IsPositive(sum))
	}

	return result.complete
}

// MultiplyInteger (MI) multiplies (U) with Aa storing the 72-bit result in Aa/Aa+1.
// Bits 0/1 of Aa/Aa+1 are sign bits.
func MultiplyInteger(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg0 := e.GetExecOrUserARegister(ci.GetA())
		aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
		res72 := pkg.Multiply(aReg0.GetW(), result.operand)
		if pkg.IsNegativeDouble(res72) {
			res72[0] |= 0_200000_000000
		}
		aReg0.SetW(res72[0])
		aReg1.SetW(res72[1])
	}

	return result.complete
}

// MultiplySingleInteger (MI) multiplies (U) with Aa storing the 36-bit result in Aa.
// If the calculation overflows 36 signed bits, DR bit 19 (overflow) is set.
// If overflow is set and operation trap is enabled, an operation trap interrupt is posted.
// How to tell from 72-bit result that overflow has occurred:
//
//	For positive numbers, word[0] and bit0 of word[1] must all be zero, else we have an overflow.
//	For negative numbers, word[0] and bit0 of word[1] must all be one, else we have an overflow.
func MultiplySingleInteger(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg := e.GetExecOrUserARegister(ci.GetA())
		res72 := pkg.Multiply(aReg.GetW(), result.operand)
		aReg.SetW(res72[1])

		okay := ((res72[0] == 0) && (res72[1]&0_400000_000000 == 0)) ||
			((res72[0] == 0_777777_777777) && (res72[1]&0_400000_000000 != 0))
		if !okay {
			dr := e.GetDesignatorRegister()
			dr.SetOverflow(true)
			if dr.IsOperationTrapEnabled() {
				i := pkg.NewOperationTrapInterrupt(pkg.OperationTrapMultiplySingleIntegerOverflow)
				e.PostInterrupt(i)
			}
		}
	}

	return result.complete
}

// MultiplyFractional (MF) multiplies (U) with Aa, performs a circular shift left by 1 bit,
// and stores the 72-bit result in Aa/Aa+1.
func MultiplyFractional(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg0 := e.GetExecOrUserARegister(ci.GetA())
		aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
		res72 := pkg.Multiply(aReg0.GetW(), result.operand)
		if pkg.IsNegativeDouble(res72) {
			res72[0] |= 0_200000_000000
		}

		res72 = pkg.LeftDoubleShiftCircular(res72, 1)
		aReg0.SetW(res72[0])
		aReg1.SetW(res72[1])
	}

	return result.complete
}

var divCheck = pkg.NewArithmeticExceptionInterrupt(pkg.ArithmeticExceptionDivideCheck)

// DivideInteger (DI) divides the 72-bit value in Aa|Aa+1 by (U),
// storing the integer quotient in Aa, and the remainder in Aa+1.
// The remainder retains the sign of the dividend.
func DivideInteger(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		divisor := result.operand

		ci := e.GetCurrentInstruction()
		aReg0 := e.GetExecOrUserARegister(ci.GetA())
		aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
		dividend := []uint64{aReg0.GetW(), aReg1.GetW()}

		quotient, remainder, divByZero, overflow := pkg.Divide(dividend, divisor)
		if divByZero || overflow {
			e.PostInterrupt(divCheck)
			return false
		}

		divIsNegative := pkg.IsNegativeDouble(dividend)
		if divIsNegative != pkg.IsNegative(divisor) {
			quotient = pkg.Negate(quotient)
		}
		if divIsNegative != pkg.IsNegative(remainder) {
			remainder = pkg.Negate(remainder)
		}

		aReg0.SetW(quotient)
		aReg1.SetW(remainder)
	}

	return result.complete
}

var signBitLookup = map[bool]uint64{
	false: 0,
	true:  pkg.NegativeZero,
}

// DivideSingleFractional (DSF) creates a dividend using Aa as the MSW and 36 sign bits as the LSW,
// right-shifts the dividend algebraically by one bit, then divides it by (U).
// The resulting quotient is stored in Aa+1, and the remainder is lost.
func DivideSingleFractional(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		divisor := result.operand

		ci := e.GetCurrentInstruction()
		aReg0 := e.GetExecOrUserARegister(ci.GetA())
		aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)

		dividend := []uint64{aReg0.GetW(), signBitLookup[pkg.IsNegative(aReg0.GetW())]}
		dividend = pkg.RightDoubleShiftAlgebraic(dividend, 1)
		if (divisor == 0) || (aReg0.GetW() > pkg.Magnitude(divisor)) {
			e.PostInterrupt(divCheck)
			return false
		}

		quotient, _, _, _ := pkg.Divide(dividend, divisor)
		if pkg.IsNegativeDouble(dividend) != pkg.IsNegative(divisor) {
			quotient = pkg.Negate(quotient)
		}

		aReg1.SetW(quotient)
	}

	return result.complete
}

// DivideFractional (DF) creates a 72-bit divisor from Aa|Aa+1, shifts it right algebraically by one bit,
// and divides it by (U) storing the integer quotient in Aa, and the remainder in Aa+1.
// The remainder retains the sign of the dividend.
func DivideFractional(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		divisor := result.operand

		ci := e.GetCurrentInstruction()
		aReg0 := e.GetExecOrUserARegister(ci.GetA())
		aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
		dividend := []uint64{aReg0.GetW(), aReg1.GetW()}
		dividend = pkg.RightDoubleShiftAlgebraic(dividend, 1)

		quotient, remainder, divByZero, overflow := pkg.Divide(dividend, divisor)
		if divByZero || overflow {
			e.PostInterrupt(divCheck)
			return false
		}

		divIsNegative := pkg.IsNegativeDouble(dividend)
		if divIsNegative != pkg.IsNegative(divisor) {
			quotient = pkg.Negate(quotient)
		}
		if divIsNegative != pkg.IsNegative(remainder) {
			remainder = pkg.Negate(remainder)
		}

		aReg0.SetW(quotient)
		aReg1.SetW(remainder)
	}

	return result.complete
}

// DoubleAddAccumulator (DA) adds (U)|(U+1) to Aa|Aa+1
func DoubleAddAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetConsecutiveOperands(true, 2, true)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg0 := e.GetExecOrUserARegister(ci.GetA())
		aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)

		addend1 := []uint64{aReg0.GetW(), aReg1.GetW()}
		addend2 := []uint64{result.source[0].GetW(), result.source[1].GetW()}
		sum := pkg.AddDouble(addend1, addend2)
		aReg0.SetW(sum[0])
		aReg1.SetW(sum[1])
		updateDesignatorRegister(e, pkg.IsPositiveDouble(addend1), pkg.IsPositiveDouble(addend2), pkg.IsPositiveDouble(sum))
	}

	return result.complete
}

// DoubleAddNegativeAccumulator (DAN) adds -(U)|(U+1) to Aa|Aa+1
func DoubleAddNegativeAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetConsecutiveOperands(true, 2, true)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg0 := e.GetExecOrUserARegister(ci.GetA())
		aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)

		addend1 := []uint64{aReg0.GetW(), aReg1.GetW()}
		addend2 := pkg.NegateDouble([]uint64{result.source[0].GetW(), result.source[1].GetW()})
		sum := pkg.AddDouble(addend1, addend2)
		fmt.Printf("add1: %012o:%012o\n", addend1[0], addend1[1]) // TODO remove
		fmt.Printf("add2: %012o:%012o\n", addend2[0], addend2[1]) // TODO remove
		fmt.Printf("sum:  %012o:%012o\n", sum[0], sum[1])         // TODO remove
		aReg0.SetW(sum[0])
		aReg1.SetW(sum[1])
		updateDesignatorRegister(e, pkg.IsPositiveDouble(addend1), pkg.IsPositiveDouble(addend2), pkg.IsPositiveDouble(sum))
	}

	return result.complete
}

//	TODO AH
//	TODO ANH
//	TODO AT
//	TODO ANT
//	TODO ADD1
//	TODO SUB1
//	TODO INC
//	TODO INC2
//	TODO DEC
//	TODO DEC2
//	TODO ENZ
//	TODO BAO
