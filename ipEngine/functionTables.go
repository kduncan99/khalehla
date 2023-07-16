// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

// FunctionTable maps the basic mode flag to either the basic mode or extended mode function table
var FunctionTable = map[bool]map[uint]func(*InstructionEngine) (completed bool, interrupt Interrupt){
	true:  BasicModeFunctionTable,
	false: ExtendedModeFunctionTable,
}

var BasicModeFunctionTable = map[uint]func(*InstructionEngine) (completed bool, interrupt Interrupt){
	010: LoadAccumulator,
	023: LoadRegister,
	027: LoadIndexRegister,
}

var ExtendedModeFunctionTable = map[uint]func(*InstructionEngine) (completed bool, interrupt Interrupt){
	010: LoadAccumulator,
	023: LoadRegister,
	027: LoadIndexRegister,
}
