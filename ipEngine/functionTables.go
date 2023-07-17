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
	001: StoreAccumulator,
	004: StoreRegister,
	005: basicModeFunction05Handler,
	006: StoreIndexRegister,
	010: LoadAccumulator,
	023: LoadRegister,
	027: LoadIndexRegister,
	074: basicModeFunction74Handler,
}

var ExtendedModeFunctionTable = map[uint]func(*InstructionEngine) (completed bool, interrupt Interrupt){
	001: StoreAccumulator,
	004: StoreRegister,
	005: extendedModeFunction05Handler,
	006: StoreIndexRegister,
	010: LoadAccumulator,
	023: LoadRegister,
	027: LoadIndexRegister,
	073: extendedModeFunction73Handler,
}

var basicModeFunction05Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt Interrupt){
	000: StoreZero,
	001: StoreNegativeZero,
	002: StorePositiveOne,
	003: StoreNegativeOne,
	004: StoreFieldataSpaces,
	005: StoreFieldataZeroes,
	006: StoreASCIISpaces,
	007: StoreASCIIZeroes,
}

var basicModeFunction74Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt Interrupt){
	006: NoOperation,
}

var extendedModeFunction05Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt Interrupt){
	000: StoreZero,
	001: StoreNegativeZero,
	002: StorePositiveOne,
	003: StoreNegativeOne,
	004: StoreFieldataSpaces,
	005: StoreFieldataZeroes,
	006: StoreASCIISpaces,
	007: StoreASCIIZeroes,
}

var extendedModeFunction73Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt Interrupt){
	014: extendedModeFunction7314Handler,
}

var extendedModeFunction7314Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt Interrupt){
	000: NoOperation,
}

func basicModeFunction05Handler(e *InstructionEngine) (completed bool, interrupt Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction05Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, NewInvalidInstructionInterrupt(InvalidInstructionBadFunctionCode)
	}
}

func basicModeFunction74Handler(e *InstructionEngine) (completed bool, interrupt Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction74Table[uint(ci.GetJ())]; found {
		return inst(e)
	} else {
		return false, NewInvalidInstructionInterrupt(InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction05Handler(e *InstructionEngine) (completed bool, interrupt Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction05Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, NewInvalidInstructionInterrupt(InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction73Handler(e *InstructionEngine) (completed bool, interrupt Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction73Table[uint(ci.GetJ())]; found {
		return inst(e)
	} else {
		return false, NewInvalidInstructionInterrupt(InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction7314Handler(e *InstructionEngine) (completed bool, interrupt Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction7314Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, NewInvalidInstructionInterrupt(InvalidInstructionBadFunctionCode)
	}
}
