// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

// FunctionTable maps the basic mode flag to either the basic mode or extended mode function table
var FunctionTable = map[bool]map[uint]func(*InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	true:  BasicModeFunctionTable,
	false: ExtendedModeFunctionTable,
}

// BasicModeFunctionTable functions indexed by the f field
var BasicModeFunctionTable = map[uint]func(*InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	001: StoreAccumulator,
	004: StoreRegister,
	005: basicModeFunction05Handler,
	006: StoreIndexRegister,
	010: LoadAccumulator,
	023: LoadRegister,
	027: LoadIndexRegister,
	073: basicModeFunction73Handler,
	074: basicModeFunction74Handler,
}

// ExtendedModeFunctionTable functions indexed by the f field
var ExtendedModeFunctionTable = map[uint]func(*InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	001: StoreAccumulator,
	004: StoreRegister,
	005: extendedModeFunction05Handler,
	006: StoreIndexRegister,
	010: LoadAccumulator,
	023: LoadRegister,
	027: LoadIndexRegister,
	073: extendedModeFunction73Handler,
}

// Basic Mode, F=005, table is indexed by the a field (most of the time the j-field indicates partial-word)
var basicModeFunction05Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: StoreZero,
	001: StoreNegativeZero,
	002: StorePositiveOne,
	003: StoreNegativeOne,
	004: StoreFieldataSpaces,
	005: StoreFieldataZeroes,
	006: StoreASCIISpaces,
	007: StoreASCIIZeroes,
}

// Basic Mode, F=073, table is indexed by the j field
var basicModeFunction73Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	015: basicModeFunction7315Handler,
}

// Basic Mode, F=073 J=015, table is indexed by the a field
var basicModeFunction7315Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	014: LoadDesignatorRegister,
	015: StoreDesignatorRegister,
}

// Basic Mode, F=074, table is indexed by the j field
var basicModeFunction74Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	006: NoOperation,
}

// Extended Mode, F=005, table is indexed by the a field (most of the time the j-field indicates partial-word)
var extendedModeFunction05Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: StoreZero,
	001: StoreNegativeZero,
	002: StorePositiveOne,
	003: StoreNegativeOne,
	004: StoreFieldataSpaces,
	005: StoreFieldataZeroes,
	006: StoreASCIISpaces,
	007: StoreASCIIZeroes,
}

// Extended Mode, F=073, table is indexed by the j field
var extendedModeFunction73Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	014: extendedModeFunction7314Handler,
	015: extendedModeFunction7315Handler,
}

// Extended Mode, F=073 J=014, table is indexed by the a field
var extendedModeFunction7314Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: NoOperation,
}

// Extended Mode, F=073 J=015, table is indexed by the a field
var extendedModeFunction7315Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	014: LoadDesignatorRegister,
	015: StoreDesignatorRegister,
}

func basicModeFunction05Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction05Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func basicModeFunction73Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction73Table[uint(ci.GetJ())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func basicModeFunction7315Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction7315Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func basicModeFunction74Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction74Table[uint(ci.GetJ())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction05Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction05Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction73Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction73Table[uint(ci.GetJ())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction7314Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction7314Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction7315Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction7315Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}
