// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

// TODO DoubleStoreAccumulator (DSA)

// StoreAccumulator (SA) stores the value of A(a) in the location indicated by U under j-field control
func StoreAccumulator(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	value := e.GetExecOrUserARegister(ci.GetA()).GetW()
	return e.StoreOperand(true, true, true, true, value)
}

// StoreASCIISpaces (SAS) stores consecutive ASCII spaces in the location indicate by U under j-field control
func StoreASCIISpaces(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	return e.StoreOperand(true, true, true, true, 0_040040_040040)
}

// StoreASCIIZeroes (SAZ) stores consecutive ASCII zeroes in the location indicate by U under j-field control
func StoreASCIIZeroes(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	return e.StoreOperand(true, true, true, true, 0_060060_060060)
}

// TODO StoreAQuarterWord (SAQW)

// StoreFieldataSpaces (SFS) stores consecutive fieldata spaces in the location indicate by U under j-field control
func StoreFieldataSpaces(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	return e.StoreOperand(true, true, true, true, 050505_050505)
}

// StoreFieldataZeroes (SFZ) stores consecutive fieldata zeroes in the location indicate by U under j-field control
func StoreFieldataZeroes(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	return e.StoreOperand(true, true, true, true, 0_606060_606060)
}

// StoreIndexRegister (SX) stores the value of X(a) in the location indicated by U under j-field control
func StoreIndexRegister(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	value := e.GetExecOrUserXRegister(ci.GetA()).GetW()
	return e.StoreOperand(true, true, true, true, value)
}

// TODO StoreMagnitudeA (SMA)
// TODO StoreNegativeA (SNA)

// StoreNegativeOne (SN1) stores a negative one in the location indicate by U under j-field control
func StoreNegativeOne(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	return e.StoreOperand(true, true, true, true, pkg.NegativeOne)
}

// StoreNegativeZero (SNZ) stores a negative zero in the location indicate by U under j-field control
func StoreNegativeZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	return e.StoreOperand(true, true, true, true, pkg.NegativeZero)
}

// StorePositiveOne (SP1) stores a positive one in the location indicate by U under j-field control
func StorePositiveOne(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	return e.StoreOperand(true, true, true, true, pkg.PositiveOne)
}

// StoreRegister (SR) stores the value of R(a) in the location indicated by U under j-field control
func StoreRegister(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	value := e.GetExecOrUserRRegister(ci.GetA()).GetW()
	return e.StoreOperand(true, true, true, true, value)
}

// TODO StoreRegisterSet (SRS)

// StoreZero (SZ) stores a positive zero in the location indicate by U under j-field control
func StoreZero(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	return e.StoreOperand(true, true, true, true, pkg.PositiveZero)
}
