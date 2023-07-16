// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/ipEngine/functions"

// FunctionTable maps the basic mode flag to either the basic mode or extended mode function table
var FunctionTable = map[bool]map[uint]func(*InstructionEngine) (completed bool, interrupt Interrupt){
	true:  BasicModeFunctionTable,
	false: ExtendedModeFunctionTable,
}

var BasicModeFunctionTable = map[uint]func(*InstructionEngine) (completed bool, interrupt Interrupt){
	001: functions.StoreAccumulator,
	004: functions.StoreRegister,
	006: functions.StoreIndexRegister,
	010: functions.LoadAccumulator,
	023: functions.LoadRegister,
	027: functions.LoadIndexRegister,
}

var ExtendedModeFunctionTable = map[uint]func(*InstructionEngine) (completed bool, interrupt Interrupt){
	001: functions.StoreAccumulator,
	004: functions.StoreRegister,
	006: functions.StoreIndexRegister,
	010: functions.LoadAccumulator,
	023: functions.LoadRegister,
	027: functions.LoadIndexRegister,
}
