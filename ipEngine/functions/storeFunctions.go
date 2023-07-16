// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package functions

import "khalehla/ipEngine"

// StoreAccumulator stores the value of A(a) in the location indicated by U under j-field control
func StoreAccumulator(e *ipEngine.InstructionEngine) (completed bool, interrupt ipEngine.Interrupt) {
	ci := e.GetCurrentInstruction()
	value := e.GetExecOrUserARegister(uint(ci.GetA())).GetW()
	return e.StoreOperand(true, true, true, true, value)
}

// StoreIndexRegister stores the value of X(a) in the location indicated by U under j-field control
func StoreIndexRegister(e *ipEngine.InstructionEngine) (completed bool, interrupt ipEngine.Interrupt) {
	ci := e.GetCurrentInstruction()
	value := e.GetExecOrUserXRegister(uint(ci.GetA())).GetW()
	return e.StoreOperand(true, true, true, true, value)
}

// StoreRegister stores the value of R(a) in the location indicated by U under j-field control
func StoreRegister(e *ipEngine.InstructionEngine) (completed bool, interrupt ipEngine.Interrupt) {
	ci := e.GetCurrentInstruction()
	value := e.GetExecOrUserRRegister(uint(ci.GetA())).GetW()
	return e.StoreOperand(true, true, true, true, value)
}
