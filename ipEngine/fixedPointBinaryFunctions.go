// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

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

//	TODO MI
//
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
	}

	return result.complete
}

//	TODO MSI
//	TODO MF
//	TODO DI
//	TODO DSF
//	TODO DF
//	TODO DA
//	TODO DAN
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
