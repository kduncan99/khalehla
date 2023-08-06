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
	011: LoadNegativeAccumulator,
	012: LoadMagnitudeAccumulator,
	013: LoadNegativeMagnitudeAccumulator,
	023: LoadRegister,
	026: LoadIndexRegisterModifier,
	027: LoadIndexRegister,
	046: LoadIndexRegisterIncrement,
	071: basicModeFunction71Handler,
	073: basicModeFunction73Handler,
	074: basicModeFunction74Handler,
	075: basicModeFunction75Handler,
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

// Basic Mode, F=071, table is indexed by the j field
var basicModeFunction71Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	013: DoubleLoadAccumulator,
	014: DoubleLoadNegativeAccumulator,
	015: DoubleLoadMagnitudeAccumulator,
}

// Basic Mode, F=073, table is indexed by the j field
var basicModeFunction73Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	015: basicModeFunction7315Handler,
	017: basicModeFunction7317Handler,
}

// Basic Mode, F=073 J=015, table is indexed by the a field
var basicModeFunction7315Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	014: LoadDesignatorRegister,
	015: StoreDesignatorRegister,
}

// Basic Mode, F=073 J=017, table is indexed by the a field
var basicModeFunction7317Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	006: InitiateAutoRecovery,
}

// Basic Mode, F=074, table is indexed by the j field
var basicModeFunction74Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	004: basicModeFunction7404Handler,
	006: NoOperation,
	015: basicModeFunction7415Handler,
}

// Basic Mode, F=074 J=04, table is indexed by the a field
var basicModeFunction7404Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: Jump,
}

// Basic Mode, F=074 J=15, table is indexed by the a field
var basicModeFunction7415Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	005: HaltJump,
}

// Basic Mode, F=075, table is indexed by the j field
var basicModeFunction75Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	013: LoadIndexRegisterLongModifier,
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

// Extended Mode, F=071, table is indexed by the j field
var extendedModeFunction71Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	013: DoubleLoadAccumulator,
	014: DoubleLoadNegativeAccumulator,
	015: DoubleLoadMagnitudeAccumulator,
}

// Extended Mode, F=073, table is indexed by the j field
var extendedModeFunction73Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	014: extendedModeFunction7314Handler,
	015: extendedModeFunction7315Handler,
	017: extendedModeFunction7317Handler,
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

// Extended Mode, F=073 J=017, table is indexed by the a field
var extendedModeFunction7317Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	006: InitiateAutoRecovery,
}

// Extended Mode, F=074, table is indexed by the j field
var extendedModeFunction74Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	015: extendedModeFunction7415Handler,
}

// Extended Mode, F=074 J=015, table is indexed by the a field
var extendedModeFunction7415Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	004: Jump,
	005: HaltJump,
}

// Extended Mode, F=075, table is indexed by the j field
var extendedModeFunction75Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	013: LoadIndexRegisterLongModifier,
}

//	Handlers -----------------------------------------------------------------------------------------------------------

func basicModeFunction05Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction05Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func basicModeFunction71Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction71Table[uint(ci.GetJ())]; found {
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

func basicModeFunction7317Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction7317Table[uint(ci.GetA())]; found {
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

func basicModeFunction7404Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction7404Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func basicModeFunction7415Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction7415Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func basicModeFunction75Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction75Table[uint(ci.GetJ())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

// ExtendedModeFunctionTable functions indexed by the f field
var ExtendedModeFunctionTable = map[uint]func(*InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	001: StoreAccumulator,
	004: StoreRegister,
	005: extendedModeFunction05Handler,
	006: StoreIndexRegister,
	010: LoadAccumulator,
	011: LoadNegativeAccumulator,
	012: LoadMagnitudeAccumulator,
	013: LoadNegativeMagnitudeAccumulator,
	023: LoadRegister,
	026: LoadIndexRegisterModifier,
	027: LoadIndexRegister,
	046: LoadIndexRegisterIncrement,
	051: LoadIndexRegisterShortIncrement,
	071: extendedModeFunction71Handler,
	073: extendedModeFunction73Handler,
	074: extendedModeFunction74Handler,
	075: extendedModeFunction75Handler,
}

func extendedModeFunction05Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction05Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction71Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction71Table[uint(ci.GetJ())]; found {
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

func extendedModeFunction7317Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction7317Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction74Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction74Table[uint(ci.GetJ())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction7415Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction7415Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction75Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction75Table[uint(ci.GetJ())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}
