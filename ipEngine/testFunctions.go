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
	completed, operand, interrupt = e.GetOperand(false, true, true, true)
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
	completed, operand, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if pkg.CountBits(operand)&01 != 0 {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

// TestLessThanEqualToModifier (TLEM)
func TestLessThanEqualToModifier(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil
	// TODO
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
	return
}

// TestNoOperation (TOP) always executes the next instruction after fetching the operand
func TestNoOperation(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil
	completed, _, interrupt = e.GetOperand(false, true, true, true)
	return
}

// TestGreaterThanZero (TGZ) skips the next instruction if U is greater than zero
func TestGreaterThanZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed = true
	interrupt = nil

	var operand uint64
	completed, operand, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if operand > 0 {
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
	completed, operand, interrupt = e.GetOperand(false, true, true, true)
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
	completed, operand, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if operand >= pkg.PositiveZero {
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
	completed, operand, interrupt = e.GetOperand(false, true, true, true)
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
	completed, operand, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if operand == pkg.NegativeZero || operand > pkg.PositiveZero {
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
	completed, operand, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if operand == pkg.PositiveZero || operand == pkg.NegativeZero {
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
	completed, operand, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if operand >= pkg.NegativeZero {
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
	completed, operand, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if operand < pkg.NegativeZero {
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
	completed, operand, interrupt = e.GetOperand(false, true, true, true)
	if completed && interrupt == nil {
		if !pkg.IsZero(operand) {
			pc := e.GetProgramAddressRegister().GetProgramCounter()
			e.SetProgramCounter(pc+2, true)
		}
	}

	return
}

//	TODO Test Positive Zero or Less Than Zero (TPZL)
//	TODO Test Not Minus Zero (TNMZ)
//	TODO Test Negative
//	TODO Test Not Positive Zero (TNPZ)
//	TODO Test Not Greater Than Zero (TNGZ)
//	TODO Test and Always Skip (TSKP)
//	TODO Test Equal (TE)
//	TODO Double-Precision Test Equal (DTE)
//	TODO Test Not Equal (TNE)
//	TODO Test Less Than or Equal (TLE)
//	TODO Test Greater (TG)
//	TODO Test Greater Magnitude (TGM)
//	TODO Double Test Greater Magnitude (DTGM)
//	TODO Test Within Range (TW)
//	TODO Test Not Within Range (TNW)
//	TODO Masked Test Equal (MTE)
//	TODO Masked Test Not Equal (MTNE)
//	TODO Masked Test Less Than or Equal (MTLE)
//	TODO Masked Test Greater (MTG)
//	TODO Masked Test Within Range (MTW)
//	TODO Masked Test Not Within Range (MTNW)
//	TODO Masked Alphanumeric Test Less Than or Equal (MATL)
//	TODO Masked Alphanumeric Test Greater (MATG)
//	TODO Test and Set (TS)
//	TODO Test and Set and Skip (TSS)
//	TODO Test and Clear and Skip (TCS)
//	TODO Conditional Replace
//	TODO Unlock (UNLK)
