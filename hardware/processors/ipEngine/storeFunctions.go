// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/common"
)

// DoubleStoreAccumulator (DSA) stores the value of A(a) and A(a+1) in the locations indicated by U and U+1
func DoubleStoreAccumulator(e *InstructionEngine) (completed bool) {
	ci := e.GetCurrentInstruction()
	aReg0 := e.GetExecOrUserARegister(ci.GetA())
	aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
	operands := []uint64{aReg0.GetW(), aReg1.GetW()}
	comp, i := e.StoreConsecutiveOperands(true, operands)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StoreAccumulator (SA) stores the value of A(a) in the location indicated by U under j-field control
func StoreAccumulator(e *InstructionEngine) (completed bool) {
	ci := e.GetCurrentInstruction()
	value := e.GetExecOrUserARegister(ci.GetA()).GetW()
	comp, i := e.StoreOperand(true, true, true, true, value)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StoreASCIISpaces (SAS) stores consecutive ASCII spaces in the location indicate by U under j-field control
func StoreASCIISpaces(e *InstructionEngine) (completed bool) {
	comp, i := e.StoreOperand(true, true, true, true, 0_040040_040040)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StoreASCIIZeroes (SAZ) stores consecutive ASCII zeroes in the location indicate by U under j-field control
func StoreASCIIZeroes(e *InstructionEngine) (completed bool) {
	comp, i := e.StoreOperand(true, true, true, true, 0_060060_060060)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StoreAQuarterWord (SAQW) stores a quarter word from register Aa into U.
// Xx.Mod is used to develop U. Xx(bit 4:5) determine which quarter word should be selected:
// value 00: Q1
// value 01: Q2
// value 02: Q3
// value 03: Q4
// The architecture leaves it undefined as to the result of setting F0.H (x-register incrementation).
// We will increment Xx in that case, which will result in strangeness, so don't set F0.H.
// It is also undefined as to what happens when F0.X is zero. We will use X0 for selecting the
// quarter-word via bits 4:5, but we will NOT use X0.Mod for developing U.
var sqwTable = []func(uint64, uint64) uint64{
	common.SetQ1,
	common.SetQ2,
	common.SetQ3,
	common.SetQ4,
}

func StoreAQuarterWord(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(false, false, false, false, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg := e.GetExecOrUserARegister(ci.GetA())
		xReg := e.GetExecOrUserXRegister(ci.GetX())

		byteSel := (xReg.GetW() >> 30) & 03
		xReg.SetW(sqwTable[byteSel](xReg.GetW(), aReg.GetW()))
	}

	return result.complete
}

// StoreFieldataSpaces (SFS) stores consecutive fieldata spaces in the location indicate by U under j-field control
func StoreFieldataSpaces(e *InstructionEngine) (completed bool) {
	comp, i := e.StoreOperand(true, true, true, true, 050505_050505)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StoreFieldataZeroes (SFZ) stores consecutive fieldata zeroes in the location indicate by U under j-field control
func StoreFieldataZeroes(e *InstructionEngine) (completed bool) {
	comp, i := e.StoreOperand(true, true, true, true, 0_606060_606060)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StoreIndexRegister (SX) stores the value of X(a) in the location indicated by U under j-field control
func StoreIndexRegister(e *InstructionEngine) (completed bool) {
	ci := e.GetCurrentInstruction()
	value := e.GetExecOrUserXRegister(ci.GetA()).GetW()

	comp, i := e.StoreOperand(true, true, true, true, value)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StoreMagnitudeA (SMA) stores the absolute value of A(a) into U
func StoreMagnitudeA(e *InstructionEngine) (completed bool) {
	ci := e.GetCurrentInstruction()
	value := e.GetExecOrUserARegister(ci.GetA()).GetW()
	if common.IsNegative(value) {
		value = common.Not(value)
	}

	comp, i := e.StoreOperand(true, true, true, true, value)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StoreNegativeA (SNA) Stores the arithmetic inverse of Aa into U
func StoreNegativeA(e *InstructionEngine) (completed bool) {
	ci := e.GetCurrentInstruction()
	value := e.GetExecOrUserARegister(ci.GetA()).GetW()
	value = common.Not(value)

	comp, i := e.StoreOperand(true, true, true, true, value)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StoreNegativeOne (SN1) stores a negative one in the location indicate by U under j-field control
func StoreNegativeOne(e *InstructionEngine) (completed bool) {
	comp, i := e.StoreOperand(true, true, true, true, common.NegativeOne)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StoreNegativeZero (SNZ) stores a negative zero in the location indicate by U under j-field control
func StoreNegativeZero(e *InstructionEngine) (completed bool) {
	comp, i := e.StoreOperand(true, true, true, true, common.NegativeZero)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StorePositiveOne (SP1) stores a positive one in the location indicate by U under j-field control
func StorePositiveOne(e *InstructionEngine) (completed bool) {
	comp, i := e.StoreOperand(true, true, true, true, common.PositiveOne)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StoreRegister (SR) stores the value of R(a) in the location indicated by U under j-field control
func StoreRegister(e *InstructionEngine) (completed bool) {
	ci := e.GetCurrentInstruction()
	value := e.GetExecOrUserRRegister(ci.GetA()).GetW()

	comp, i := e.StoreOperand(true, true, true, true, value)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// StoreRegisterSet (LRS) Stores the GRS (or one or two subsets thereof) into U through U+n.
// Specifically, the instruction defines two sets of ranges and lengths as follows:
// Aa[2:8]   = range 2 length
// Aa[11:17] = range 2 first GRS index
// Aa[20:26] = range 1 count
// Aa[29:35] = range 1 first GRS index
// So we start storing registers from GRS index of range 1, for the number of registers in range 1 count,
// to U[0] to U[range1count - 1], and then from GRS index of range 2, for the number of registers in range 2 count,
// to U[range1count] to U[range1count + range2count - 1].
// If either count is zero, then the associated range is not used.
// If the GRS address exceeds 0177, it wraps around to zero.
func StoreRegisterSet(e *InstructionEngine) (completed bool) {
	ci := e.GetCurrentInstruction()
	aReg := e.GetExecOrUserARegister(ci.GetA())

	count2 := aReg.GetQ1() & 0177
	address2 := aReg.GetQ2() & 0177
	count1 := aReg.GetQ3() & 0177
	address1 := aReg.GetQ4() & 0177

	result := e.GetConsecutiveOperands(false, count1+count2, true)
	if result.complete && result.interrupt == nil {
		dr := e.GetDesignatorRegister()
		grs := e.GetGeneralRegisterSet()
		opx := 0

		if count1 > 0 {
			grsIndex := address1
			for x := 0; x < int(count1); x++ {
				if !e.isGRSAccessAllowed(grsIndex, dr.GetProcessorPrivilege(), true) {
					i := common.NewReferenceViolationInterrupt(common.ReferenceViolationReadAccess, false)
					e.PostInterrupt(i)
					return false
				}

				grs.SetRegisterValue(grsIndex, result.source[opx].GetW())

				opx++
				grsIndex++
				if grsIndex == 0200 {
					grsIndex = 0
				}
			}
		}

		if count2 > 0 {
			grsIndex := address2
			for x := 0; x < int(count2); x++ {
				if !e.isGRSAccessAllowed(grsIndex, dr.GetProcessorPrivilege(), true) {
					i := common.NewReferenceViolationInterrupt(common.ReferenceViolationReadAccess, false)
					e.PostInterrupt(i)
					return false
				}

				grs.SetRegisterValue(grsIndex, result.source[opx].GetW())

				opx++
				grsIndex++
				if grsIndex == 0200 {
					grsIndex = 0
				}
			}
		}
	}

	return true
}

// StoreZero (SZ) stores a positive zero in the location indicate by U under j-field control
func StoreZero(e *InstructionEngine) (completed bool) {
	comp, i := e.StoreOperand(true, true, true, true, common.PositiveZero)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}
