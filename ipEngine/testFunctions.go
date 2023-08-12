// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

// TestEvenParity (TEP) skips the next instruction if U has an even number of bits set to one.
func TestEvenParity(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if pkg.CountBits(operand)&01 == 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestOddParity (TOP) skips the next instruction if U has an odd number of bits set to one.
func TestOddParity(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if pkg.CountBits(operand)&01 != 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestLessThanOrEqualToModifier (TLEM)
// This one is a little more complicated. At heart, we're just dealing with XI and XM for Xa.
// The comparison is between the unsigned lower 18 bits of (U), and the unsigned XM portion of Xa.
// After the comparison, XI is added to XM of Xa.
// Then we evaluate the comparison; if it passes, we skip the next instruction.
// Things get a bit weird when U is developed with an index register Xx where Xa is Xx...
// In basic mode, when Xx os Xa and auto-increment is set on the instruction, we do *not*
// increment the index register twice (as one might think we would)... the increment only happens
// once.
// 24-bit indexing does not apply to the comparison or increment of Xa, but it *does* (still) apply
// to developing U (if it were to apply at all). There is no conflict with the previous point,
// as 24-bit indexing can only be in effect during extended mode.
func TestLessThanOrEqualToModifier(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil
	// TODO remove these comments
	// The contents of U are fetched under F0.j control, its high order 18 bits are truncated, and it is
	// alphanumerically compared with the right half (bits 18–35) of Xa;
	// bits 0–17 of Xa are added to bits 18–35 of Xa and the sum is stored into bits 18–35 of Xa.
	// Bits 0–17 of Xa remain unchanged.
	// If the results of the comparison indicate that the operand is less than or equal to Xa, the next
	// instruction (NI) is skipped and the instruction following NI is the next executed; otherwise, NI is
	// the next instruction executed.

	// If F0.a = 0, Index_Register X0 is referenced.
	// • Both Xa bits 18–35 and the value from U are considered to be 18-bit values, and they are
	//   compared alphanumerically, where the sign bit is treated as a data bit (hence –0 > +0).
	// • As only the rightmost 18 bits of the value from U are involved in the operation, F0.j = 0,
	//   F0.j = 1, or F0.j = 3 yield the same results. F0.j = 016 or F0.j = 017 yield the same results.
	// • In Basic_Mode, if F0.h = 1 and F0.a =F0. x, the specified X-Register is incremented or
	//   modified only once.
	// • 24-bit indexing applies only to the fetch of the instruction operand and not to the operand
	//   comparison or addition.

	//	Grab the value of Xa MOD first, before address resolution has a chance to auto-increment
	//	Xx (which might be the same X register as Xa). We want to compare against the modifier
	//	as it was *before* auto-incrementation.
	ci := e.GetCurrentInstruction()
	xa := e.GetExecOrUserXRegister(ci.GetA())
	xaMod := xa.GetXM()

	//	Now develop U.
	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if !completed || interrupt != nil {
		return
	}

	//	perform the comparison
	if (operand & 0_777777) <= xaMod {
		if pkg.CountBits(operand)&01 != 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	//	now (probably) increment Xa.
	basic := e.activityStatePacket.GetDesignatorRegister().IsBasicModeEnabled()
	sameXReg := ci.GetA() == ci.GetX()
	incrXReg := ci.GetH() == 1
	if !(basic && sameXReg && incrXReg) {
		xa.IncrementModifier()
	}

	return
}

// TestNoOperation (TOP) always executes the next instruction after fetching the operand
func TestNoOperation(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil
	_, completed, interrupt = e.GetOperand(false, true, true, true)
	return
}

// TestGreaterThanZero (TGZ) skips the next instruction if U is greater than zero
func TestGreaterThanZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if operand != pkg.PositiveZero && pkg.IsPositive(operand) {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestPositiveZero (TPZ) skips the next instruction if U is equal to positive zero
func TestPositiveZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if operand == pkg.PositiveZero {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestPositive (TPZ) skips the next instruction if U is greater than or equal to positive zero
func TestPositive(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if pkg.IsPositive(operand) {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestMinusZero (TMZ) skips the next instruction if U is equal to negative zero
func TestMinusZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if operand == pkg.NegativeZero {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestMinusZeroOrGreaterThanZero (TMZG) skips the next instruction if U is equal to negative zero
// or greater than positive zero.
func TestMinusZeroOrGreaterThanZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if operand == pkg.NegativeZero ||
			(operand != pkg.PositiveZero && pkg.IsPositive(operand)) {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestZero (TZ) skips the next instruction if U is positive or negative zero
func TestZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if pkg.IsZero(operand) {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestNotLessThanZero (TNLZ) skips the next instruction if U is not less than negative zero
func TestNotLessThanZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if pkg.IsPositive(operand) && operand != pkg.PositiveZero {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestLessThanZero (TNLZ) skips the next instruction if U is less than negative zero
func TestLessThanZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if pkg.IsNegative(operand) && operand != pkg.NegativeZero {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestNonZero skips the next instruction if U is *not* positive or negative zero
func TestNonZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if !pkg.IsZero(operand) {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestPositiveZeroOrLessThanZero skips the next instruction if U is positive zero, or less than negative zero
func TestPositiveZeroOrLessThanZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if (operand == pkg.PositiveZero) ||
			(pkg.IsNegative(operand) && operand != pkg.NegativeZero) {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestNotMinusZero skips the next instruction if U is *not* negative zero
func TestNotMinusZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if operand != pkg.NegativeZero {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestNegative skips the next instruction if U is negative
func TestNegative(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if !pkg.IsNegative(operand) {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestNotPositiveZero skips the next instruction if U is *not* positive zero
func TestNotPositiveZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if operand != pkg.PositiveZero {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestNotGreaterThanZero skips the next instruction if U is not greater than (positive) zero
func TestNotGreaterThanZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if pkg.IsNegative(operand) || operand == pkg.PositiveZero {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestAndAlwaysSkip always skips the next instruction, ignoring the operand
func TestAndAlwaysSkip(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	_, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		pc := e.GetProgramAddressRegister().GetProgramCounter()
		e.SetProgramCounter(pc+2, true)
	}

	return
}

// TestEqual (TE) skips the next instruction if the operand is equal to Aa.
// Positive zero is not equal to negative zero.
func TestEqual(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserARegister(ci.GetA()).GetW()
		if operand == aValue {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// DoubleTestEqual (TE) skips the next instruction if the 72-bit operand is equal to Aa | Aa+1
// Positive zero is not equal to negative zero.
func DoubleTestEqual(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operands []pkg.Word36
	operands, completed, interrupt = e.GetConsecutiveOperands(true, 2, false)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		ax := e.GetExecOrUserARegisterIndex(ci.GetA())
		aRegs := e.generalRegisterSet.registers[ax : ax+2]
		if aRegs[0] == operands[0] && aRegs[1] == operands[1] {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestNotEqual (TNE) skips the next instruction if the operand is not equal to Aa
// Positive zero is not equal to negative zero.
func TestNotEqual(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserARegister(ci.GetA()).GetW()
		if operand != aValue {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestLessThanOrEqual (TLE) skips the next instruction if the operand is less than or equal to Aa
// Positive zero is greater than negative zero.
func TestLessThanOrEqual(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserARegister(ci.GetA()).GetW()
		if pkg.Compare(operand, aValue) <= 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestGreater (TG) skips the next instruction if the operand is greater than Aa
// Positive zero is greater than negative zero.
func TestGreater(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserARegister(ci.GetA()).GetW()
		if pkg.Compare(operand, aValue) > 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestGreaterMagnitude (TGM) skips the next instruction if the magnitude of the operand is greater than Aa
func TestGreaterMagnitude(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserARegister(ci.GetA()).GetW()
		if pkg.Compare(pkg.Magnitude(operand), aValue) > 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// DoubleTestGreaterMagnitude (DTGM) skips the next instruction if the magnitude of U|U+1
// is greater than the value of Aa|Aa+1.
func DoubleTestGreaterMagnitude(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserARegister(ci.GetA()).GetW()
		if pkg.Compare(pkg.Magnitude(operand), aValue) > 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestWithinRange (TW) skips the next instruction if the operand is greater than Aa and less than or equal to Aa+1
func TestWithinRange(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		ax := e.GetExecOrUserARegisterIndex(ci.GetA())
		a1 := e.GetExecOrUserARegister(ax).GetW()
		a2 := e.GetExecOrUserARegister(ax + 1).GetW()
		if pkg.Compare(pkg.Magnitude(operand), a1) > 0 &&
			pkg.Compare(pkg.Magnitude(operand), a2) <= 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestNotWithinRange (TW) skips the next instruction if the operand is not greater than Aa
// or not less than or equal to Aa+1
func TestNotWithinRange(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		ax := e.GetExecOrUserARegisterIndex(ci.GetA())
		a1 := e.GetExecOrUserARegister(ax).GetW()
		a2 := e.GetExecOrUserARegister(ax + 1).GetW()
		if pkg.Compare(pkg.Magnitude(operand), a1) <= 0 &&
			pkg.Compare(pkg.Magnitude(operand), a2) > 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// MaskedTestEqual (MTE) skips the next instruction if the operand AND R2 is equal to Aa AND R2
func MaskedTestEqual(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, false, false)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserARegister(ci.GetA()).GetW()
		rValue := e.GetExecOrUserRRegister(R2).GetW()
		if pkg.And(operand, rValue) == pkg.And(aValue, rValue) {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// MaskedTestNotEqual (MTNE) skips the next instruction if the operand AND R2 is *not* equal to Aa AND R2
func MaskedTestNotEqual(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, false, false)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserARegister(ci.GetA()).GetW()
		rValue := e.GetExecOrUserRRegister(R2).GetW()
		if pkg.And(operand, rValue) != pkg.And(aValue, rValue) {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// MaskedTestLessThanOrEqual (MTLE) skips the next instruction
// if the operand AND R2 is less than or equal to Aa AND R2
func MaskedTestLessThanOrEqual(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, false, false)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserARegister(ci.GetA()).GetW()
		rValue := e.GetExecOrUserRRegister(R2).GetW()
		if pkg.Compare(pkg.And(operand, rValue), pkg.And(aValue, rValue)) <= 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// MaskedTestGreater (MTG) skips the next instruction
// if the operand AND R2 is greater than Aa AND R2
func MaskedTestGreater(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, false, false)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserARegister(ci.GetA()).GetW()
		rValue := e.GetExecOrUserRRegister(R2).GetW()
		if pkg.Compare(pkg.And(operand, rValue), pkg.And(aValue, rValue)) > 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// MaskedTestWithinRange (MTW) skips the next instruction if the operand AND R2
// is greater than Aa AND R2 and less than or equal to Aa+1 AND R2
func MaskedTestWithinRange(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		ax := e.GetExecOrUserARegisterIndex(ci.GetA())
		rVal := e.GetExecOrUserRRegister(R2).GetW()
		a1Masked := pkg.And(e.GetExecOrUserARegister(ax).GetW(), rVal)
		a2Masked := pkg.And(e.GetExecOrUserARegister(ax+1).GetW(), rVal)
		opMasked := pkg.And(operand, rVal)

		if pkg.Compare(pkg.Magnitude(opMasked), a1Masked) > 0 &&
			pkg.Compare(pkg.Magnitude(opMasked), a2Masked) <= 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// MaskedTestNotWithinRange (MTNW) skips the next instruction if the operand AND R2
// is not greater than Aa AND R2 or not less than or equal to Aa+1 AND R2
func MaskedTestNotWithinRange(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		ax := e.GetExecOrUserARegisterIndex(ci.GetA())
		rVal := e.GetExecOrUserRRegister(R2).GetW()
		a1Masked := pkg.And(e.GetExecOrUserARegister(ax).GetW(), rVal)
		a2Masked := pkg.And(e.GetExecOrUserARegister(ax+1).GetW(), rVal)
		opMasked := pkg.And(operand, rVal)

		if pkg.Compare(pkg.Magnitude(opMasked), a1Masked) <= 0 ||
			pkg.Compare(pkg.Magnitude(opMasked), a2Masked) > 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// MaskedAlphanumericTestLessThanOrEqual (MATL) skips the next instruction
// if the operand AND R2 is alphanumerically less than or equal to Aa AND R2
// For this test, alphanumeric simply means that bit 0 is data, not a sign bit,
// so we can do a simple binary comparison.
func MaskedAlphanumericTestLessThanOrEqual(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, false, false)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserARegister(ci.GetA()).GetW()
		rValue := e.GetExecOrUserRRegister(R2).GetW()
		if pkg.And(operand, rValue) <= pkg.And(aValue, rValue) {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// MaskedAlphanumericTestGreater (MTG) skips the next instruction
// if the operand AND R2 is greater than Aa AND R2
// For this test, alphanumeric simply means that bit 0 is data, not a sign bit,
// so we can do a simple binary comparison.
func MaskedAlphanumericTestGreater(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	operand, completed, interrupt = e.GetOperand(false, true, false, false)
	if completed && interrupt == nil {
		ci := e.GetCurrentInstruction()
		aValue := e.GetExecOrUserARegister(ci.GetA()).GetW()
		rValue := e.GetExecOrUserRRegister(R2).GetW()
		if pkg.And(operand, rValue) > pkg.And(aValue, rValue) {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestAndSet (TS)
func TestAndSet(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	//	TODO

	return
}

// TestAndSetAndSkip (TSS)
func TestAndSetAndSkip(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	//	TODO

	return
}

// TestAndClearAndSkip (TCS)
func TestAndClearAndSkip(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	//	TODO

	return
}

// ConditionalReplace (CR)
func ConditionalReplace(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	//	TODO

	return
}

// Unlock (UNLK)
func Unlock(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	//	TODO

	return
}
