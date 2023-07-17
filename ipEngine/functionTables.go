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
	005: basicModeFunction05Handler,
	006: functions.StoreIndexRegister,
	010: functions.LoadAccumulator,
	023: functions.LoadRegister,
	027: functions.LoadIndexRegister,
	074: basicModeFunction74Handler,
}

var ExtendedModeFunctionTable = map[uint]func(*InstructionEngine) (completed bool, interrupt Interrupt){
	001: functions.StoreAccumulator,
	004: functions.StoreRegister,
	005: extendedModeFunction05Handler,
	006: functions.StoreIndexRegister,
	010: functions.LoadAccumulator,
	023: functions.LoadRegister,
	027: functions.LoadIndexRegister,
	073: extendedModeFunction73Handler,
}

var basicModeFunction05Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt Interrupt){
	000: functions.StoreZero,
	001: functions.StoreNegativeZero,
	002: functions.StorePositiveOne,
	003: functions.StoreNegativeOne,
	004: functions.StoreFieldataSpaces,
	005: functions.StoreFieldataZeroes,
	006: functions.StoreASCIISpaces,
	007: functions.StoreASCIIZeroes,
}

var basicModeFunction74Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt Interrupt){
	006: functions.NoOperation,
}

var extendedModeFunction05Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt Interrupt){
	000: functions.StoreZero,
	001: functions.StoreNegativeZero,
	002: functions.StorePositiveOne,
	003: functions.StoreNegativeOne,
	004: functions.StoreFieldataSpaces,
	005: functions.StoreFieldataZeroes,
	006: functions.StoreASCIISpaces,
	007: functions.StoreASCIIZeroes,
}

var extendedModeFunction73Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt Interrupt){
	014: extendedModeFunction7314Handler,
}

var extendedModeFunction7314Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt Interrupt){
	000: functions.NoOperation,
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
