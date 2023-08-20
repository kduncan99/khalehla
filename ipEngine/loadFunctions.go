// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

// DoubleLoadAccumulator (DL) loads the content of U and U+1, storing the values in Aa and Aa+1
func DoubleLoadAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetConsecutiveOperands(true, 2, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		grsIndex := e.GetExecOrUserARegisterIndex(ci.GetA())
		e.generalRegisterSet.SetRegisterValue(grsIndex, result.source[0].GetW())
		e.generalRegisterSet.SetRegisterValue(grsIndex+1, result.source[1].GetW())
	}

	return result.complete
}

// DoubleLoadMagnitudeAccumulator (DL) loads the arithmetic magnitude of the content of U and U+1,
// storing the values in Aa and Aa+1
func DoubleLoadMagnitudeAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetConsecutiveOperands(true, 2, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		grsIndex := e.GetExecOrUserARegisterIndex(ci.GetA())
		if result.source[0].IsNegative() {
			e.generalRegisterSet.SetRegisterValue(grsIndex, pkg.Not(result.source[0].GetW()))
			e.generalRegisterSet.SetRegisterValue(grsIndex+1, pkg.Not(result.source[1].GetW()))
		} else {
			e.generalRegisterSet.SetRegisterValue(grsIndex, result.source[0].GetW())
			e.generalRegisterSet.SetRegisterValue(grsIndex+1, result.source[1].GetW())
		}
	}

	return result.complete
}

// DoubleLoadNegativeAccumulator (DL) loads the arithmetic negative of the content of U and U+1,
// storing the values in Aa and Aa+1
func DoubleLoadNegativeAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetConsecutiveOperands(true, 2, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		grsIndex := e.GetExecOrUserARegisterIndex(ci.GetA())
		e.generalRegisterSet.SetRegisterValue(grsIndex, pkg.Not(result.source[0].GetW()))
		e.generalRegisterSet.SetRegisterValue(grsIndex+1, pkg.Not(result.source[1].GetW()))
	}

	return result.complete
}

// LoadAccumulator (LA) loads the content of U under j-field control, and stores it in A(a)
func LoadAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(result.operand)
	}

	return result.complete
}

// LoadAQuarterWord (LAQW) loads a quarter word from U into register Aa.
// Xx.Mod is used to develop U. Xx(bit 4:5) determine which quarter word should be selected:
// value 00: Q1
// value 01: Q2
// value 02: Q3
// value 03: Q4
// The architecture leaves it undefined as to the result of setting F0.H (x-register incrementation).
// We will increment Xx in that case, which will result in strangeness, so don't set F0.H.
// It is also undefined as to what happens when F0.X is zero. We will use X0 for selecting the
// quarter-word via bits 4:5, but we will NOT use X0.Mod for developing U.
var lqwTable = []func(uint64) uint64{
	pkg.GetQ1,
	pkg.GetQ2,
	pkg.GetQ3,
	pkg.GetQ4,
}

func LoadAQuarterWord(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(false, false, false, false, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		aReg := e.GetExecOrUserARegister(ci.GetA())
		xReg := e.GetExecOrUserXRegister(ci.GetX())

		byteSel := (xReg.GetW() >> 30) & 03
		value := lqwTable[byteSel](result.operand)
		aReg.SetW(value)
	}

	return result.complete
}

// LoadIndexRegister (LX) loads the content of U under j-field control, and stores it in X(a)
func LoadIndexRegister(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetW(result.operand)
	}

	return result.complete
}

// LoadIndexRegisterModifier (LXM)
func LoadIndexRegisterModifier(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetXM(result.operand)
	}

	return result.complete
}

// LoadIndexRegisterLongModifier (LXLM)
func LoadIndexRegisterLongModifier(e *InstructionEngine) (completed bool) {
	dr := e.activityStatePacket.GetDesignatorRegister()
	if dr.IsBasicModeEnabled() && dr.GetProcessorPrivilege() > 0 {
		i := pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
		e.PostInterrupt(i)
		return false
	}

	result := e.GetOperand(true, true, false, false, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetXM24(result.operand)
	}

	return result.complete
}

// LoadIndexRegisterIncrement (LXI)
func LoadIndexRegisterIncrement(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetXI(result.operand)
	}

	return result.complete
}

// LoadIndexRegisterShortIncrement (LXSI)
func LoadIndexRegisterShortIncrement(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserXRegister(ci.GetA()).SetXI12(result.operand)
	}

	return result.complete
}

// LoadMagnitudeAccumulator (LMA)
func LoadMagnitudeAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		if pkg.IsNegative(result.operand) {
			result.operand ^= pkg.NegativeZero
		}
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(result.operand)
	}

	return result.complete
}

// LoadNegativeAccumulator (LNA)
func LoadNegativeAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		result.operand ^= pkg.NegativeZero
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(result.operand)
	}

	return result.complete
}

// LoadNegativeMagnitudeAccumulator (LNMA)
func LoadNegativeMagnitudeAccumulator(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		if !pkg.IsNegative(result.operand) {
			result.operand ^= pkg.NegativeZero
		}
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(result.operand)
	}

	return result.complete
}

// LoadRegister (LR) loads the content of U under j-field control, and stores it in R(a)
func LoadRegister(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserRRegister(ci.GetA()).SetW(result.operand)
	}

	return result.complete
}

// LoadRegisterSet (LRS) Loads the GRS (or one or two subsets thereof) from the contents of U through U+n.
// Specifically, the instruction defines two sets of ranges and lengths as follows:
// Aa[2:8]   = range 2 length
// Aa[11:17] = range 2 first GRS index
// Aa[20:26] = range 1 count
// Aa[29:35] = range 1 first GRS index
// So we start loading registers from GRS index of range 1, for the number of registers in range 1 count,
// from U[0] to U[range1count - 1], and then from GRS index of range 2, for the number of registers in range 2 count,
// from U[range1count] to U[range1count + range2count - 1].
// If either count is zero, then the associated range is not used.
// If the GRS address exceeds 0177, it wraps around to zero.
func LoadRegisterSet(e *InstructionEngine) (completed bool) {
	ci := e.GetCurrentInstruction()
	aReg := e.GetExecOrUserARegister(ci.GetA())
	count2 := aReg.GetQ1() & 0177
	address2 := aReg.GetQ2() & 0177
	count1 := aReg.GetQ3() & 0177
	address1 := aReg.GetQ4() & 0177

	result := e.GetConsecutiveOperands(false, count1+count2, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		grs := e.GetGeneralRegisterSet()
		dr := e.GetDesignatorRegister()

		ux := 0
		grsIndex := address1
		count := count1
		for count > 0 {
			if !e.isGRSAccessAllowed(grsIndex, dr.GetProcessorPrivilege(), true) {
				i := pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationReadAccess, false)
				e.PostInterrupt(i)
				return false
			}

			grs.registers[grsIndex].SetW(result.source[ux].GetW())
			ux++
			grsIndex++
			if grsIndex == 0200 {
				grsIndex = 0
			}
			count--
		}

		grsIndex = address2
		count = count2
		for count > 0 {
			if !e.isGRSAccessAllowed(grsIndex, dr.GetProcessorPrivilege(), true) {
				i := pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationReadAccess, false)
				e.PostInterrupt(i)
				return false
			}

			grs.registers[grsIndex].SetW(result.source[ux].GetW())
			ux++
			grsIndex++
			if grsIndex == 0200 {
				grsIndex = 0
			}
			count--
		}
	}

	return result.complete
}

// LoadStringBitLength (LSBL) Copies the right-most 6 bits of U to Xa bits 6-11
func LoadStringBitLength(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		xReg := e.GetExecOrUserXRegister(ci.GetA())
		value := (xReg.GetW() & 0_770077_777777) | ((result.operand & 077) << 24)
		xReg.SetW(value)
	}

	return result.complete
}

// LoadStringBitOffset (LSBO) Copies the right-most 6 bits of U to Xa bits 0-5
func LoadStringBitOffset(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(true, true, true, true, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
		return false
	} else if result.complete {
		ci := e.GetCurrentInstruction()
		xReg := e.GetExecOrUserXRegister(ci.GetA())
		value := (xReg.GetW() & 0_007777_777777) | ((result.operand & 077) << 30)
		xReg.SetW(value)
	}

	return result.complete
}
