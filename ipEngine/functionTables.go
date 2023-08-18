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
	002: StoreNegativeA,
	003: StoreMagnitudeA,
	004: StoreRegister,
	005: basicModeFunction05Handler,
	006: StoreIndexRegister,
	007: basicModeFunction07Handler,
	010: LoadAccumulator,
	011: LoadNegativeAccumulator,
	012: LoadMagnitudeAccumulator,
	013: LoadNegativeMagnitudeAccumulator,
	023: LoadRegister,
	026: LoadIndexRegisterModifier,
	027: LoadIndexRegister,
	044: TestEvenParity,
	045: TestOddParity,
	046: LoadIndexRegisterIncrement,
	047: TestLessThanOrEqualToModifier,
	050: TestZero,
	051: TestNonZero,
	052: TestEqual,
	053: TestNotEqual,
	054: TestLessThanOrEqual,
	055: TestGreater,
	056: TestWithinRange,
	057: TestNotWithinRange,
	060: TestPositive,
	061: TestNegative,
	070: JumpGreaterAndDecrement,
	071: basicModeFunction71Handler,
	072: basicModeFunction72Handler,
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

// Basic Mode, F=007, table is indexed by the j field
var basicModeFunction07Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	004: LoadAQuarterWord,
	005: StoreAQuarterWord,
	014: LoadProgramControlDesignators,
	015: StoreProgramControlDesignators,
}

// Basic Mode, F=071, table is indexed by the j field
var basicModeFunction71Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	012: DoubleStoreAccumulator,
	013: DoubleLoadAccumulator,
	014: DoubleLoadNegativeAccumulator,
	015: DoubleLoadMagnitudeAccumulator,
	016: DoubleJumpZero,
	017: DoubleTestEqual,
}

// Basic Mode, F=072, table is indexed by the j field
var basicModeFunction72Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	001: StoreLocationAndJump,
	002: JumpPositiveAndShift,
	003: JumpNegativeAndShift,
	016: StoreRegisterSet,
	017: LoadRegisterSet,
}

// Basic Mode, F=073, table is indexed by the j field
var basicModeFunction73Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	015: basicModeFunction7315Handler,
	017: basicModeFunction7317Handler,
}

// Basic Mode, F=073 J=015, table is indexed by the a field
var basicModeFunction7315Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	003: AccelerateUserRegisterSet,
	004: DecelerateUserRegisterSet,
	014: LoadDesignatorRegister,
	015: StoreDesignatorRegister,
}

// Basic Mode, F=073 J=017, table is indexed by the a field
var basicModeFunction7317Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: TestAndSet,
	001: TestAndSetAndSkip,
	002: TestAndClearAndSkip,
	006: InitiateAutoRecovery,
}

// Basic Mode, F=074, table is indexed by the j field
var basicModeFunction74Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: JumpZero,
	001: JumpNonZero,
	002: JumpPositive,
	003: JumpNegative,
	004: basicModeFunction7404Handler,
	005: HaltKeysAndJump,
	006: NoOperation,
	010: JumpNoLowBit,
	011: JumpLowBit,
	012: JumpModifierGreaterAndIncrement,
	013: LoadModifierAndJump,
	014: basicModeFunction7414Handler,
	015: basicModeFunction7415Handler,
	016: JumpCarry,
	017: JumpNoCarry,
}

// Basic Mode, F=074 J=04, table is indexed by the a field
var basicModeFunction7404Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: Jump,
	001: JumpKeys,
	002: JumpKeys,
	003: JumpKeys,
	004: JumpKeys,
	005: JumpKeys,
	006: JumpKeys,
	007: JumpKeys,
	010: JumpKeys,
	011: JumpKeys,
	012: JumpKeys,
	013: JumpKeys,
	014: JumpKeys,
	015: JumpKeys,
	016: JumpKeys,
	017: JumpKeys,
}

// Basic Mode, F=074 J=14, table is indexed by the a field
var basicModeFunction7414Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: JumpOverflow,
	001: JumpFloatingUnderflow,
	002: JumpFloatingOverflow,
	003: JumpDivideFault,
}

// Basic Mode, F=074 J=15, table is indexed by the a field
var basicModeFunction7415Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: JumpNoOverflow,
	001: JumpNoFloatingUnderflow,
	002: JumpNoFloatingOverflow,
	003: JumpNoDivideFault,
	005: HaltJump,
}

// Basic Mode, F=075, table is indexed by the j field
var basicModeFunction75Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	013: LoadIndexRegisterLongModifier,
	015: ConditionalReplace,
}

//	--------------------------------------------------------------------------------------------------------------------

// ExtendedModeFunctionTable functions indexed by the f field
var ExtendedModeFunctionTable = map[uint]func(*InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	001: StoreAccumulator,
	002: StoreNegativeA,
	003: StoreMagnitudeA,
	004: StoreRegister,
	005: extendedModeFunction05Handler,
	006: StoreIndexRegister,
	007: extendedModeFunction07Handler,
	010: LoadAccumulator,
	011: LoadNegativeAccumulator,
	012: LoadMagnitudeAccumulator,
	013: LoadNegativeMagnitudeAccumulator,
	023: LoadRegister,
	026: LoadIndexRegisterModifier,
	027: LoadIndexRegister,
	033: extendedModeFunction33Handler,
	044: TestEvenParity,
	045: TestOddParity,
	046: LoadIndexRegisterIncrement,
	047: TestLessThanOrEqualToModifier,
	050: extendedModeFunction50Handler,
	051: LoadIndexRegisterShortIncrement,
	052: TestEqual,
	053: TestNotEqual,
	054: TestLessThanOrEqual,
	055: TestGreater,
	056: TestWithinRange,
	057: TestNotWithinRange,
	060: LoadStringBitOffset,
	061: LoadStringBitLength,
	070: JumpGreaterAndDecrement,
	071: extendedModeFunction71Handler,
	072: extendedModeFunction72Handler,
	073: extendedModeFunction73Handler,
	074: extendedModeFunction74Handler,
	075: extendedModeFunction75Handler,
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

// Extended Mode, F=007, table is indexed by the j field
var extendedModeFunction07Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	004: LoadAQuarterWord,
	005: StoreAQuarterWord,
}

// Extended Mode, F=033, table is indexed by the j field
var extendedModeFunction33Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	013: TestGreaterMagnitude,
	014: DoubleTestGreaterMagnitude,
}

// Extended Mode, F=050, table is indexed by the a field
var extendedModeFunction50Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: TestNoOperation,
	001: TestGreaterThanZero,
	002: TestPositiveZero,
	003: TestPositive,
	004: TestMinusZero,
	005: TestMinusZeroOrGreaterThanZero,
	006: TestZero,
	007: TestNotLessThanZero,
	010: TestLessThanZero,
	011: TestNonZero,
	012: TestPositiveZeroOrLessThanZero,
	013: TestNotMinusZero,
	014: TestNegative,
	015: TestNotPositiveZero,
	016: TestNotGreaterThanZero,
	017: TestAndAlwaysSkip,
}

// Extended Mode, F=071, table is indexed by the j field
var extendedModeFunction71Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: MaskedTestEqual,
	001: MaskedTestNotEqual,
	002: MaskedTestLessThanOrEqual,
	003: MaskedTestGreater,
	004: MaskedTestWithinRange,
	005: MaskedTestNotWithinRange,
	006: MaskedAlphanumericTestLessThanOrEqual,
	007: MaskedAlphanumericTestGreater,
	012: DoubleStoreAccumulator,
	013: DoubleLoadAccumulator,
	014: DoubleLoadNegativeAccumulator,
	015: DoubleLoadMagnitudeAccumulator,
	016: DoubleJumpZero,
	017: DoubleTestEqual,
}

// Extended Mode, F=072, table is indexed by the j field
var extendedModeFunction72Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	002: JumpPositiveAndShift,
	003: JumpNegativeAndShift,
	016: StoreRegisterSet,
	017: LoadRegisterSet,
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
	004: Unlock,
}

// Extended Mode, F=073 J=015, table is indexed by the a field
var extendedModeFunction7315Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	003: AccelerateUserRegisterSet,
	004: DecelerateUserRegisterSet,
	014: LoadDesignatorRegister,
	015: StoreDesignatorRegister,
}

// Extended Mode, F=073 J=017, table is indexed by the a field
var extendedModeFunction7317Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: TestAndSet,
	001: TestAndSetAndSkip,
	002: TestAndClearAndSkip,
	004: LoadUserDesignators,
	005: StoreUserDesignators,
	006: InitiateAutoRecovery,
}

// Extended Mode, F=074, table is indexed by the j field
var extendedModeFunction74Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: JumpZero,
	001: JumpNonZero,
	002: JumpPositive,
	003: JumpNegative,
	010: JumpNoLowBit,
	011: JumpLowBit,
	012: JumpModifierGreaterAndIncrement,
	013: LoadModifierAndJump,
	014: extendedModeFunction7414Handler,
	015: extendedModeFunction7415Handler,
}

// Extended Mode, F=074 J=014, table is indexed by the a field
var extendedModeFunction7414Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: JumpOverflow,
	001: JumpFloatingUnderflow,
	002: JumpFloatingOverflow,
	003: JumpDivideFault,
	004: JumpCarry,
	005: JumpNoCarry,
}

// Extended Mode, F=074 J=015, table is indexed by the a field
var extendedModeFunction7415Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	000: JumpNoOverflow,
	001: JumpNoFloatingUnderflow,
	002: JumpNoFloatingOverflow,
	003: JumpNoDivideFault,
	004: Jump,
	005: HaltJump,
}

// Extended Mode, F=075, table is indexed by the j field
var extendedModeFunction75Table = map[uint]func(engine *InstructionEngine) (completed bool, interrupt pkg.Interrupt){
	013: LoadIndexRegisterLongModifier,
	015: ConditionalReplace,
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

func basicModeFunction07Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction07Table[uint(ci.GetJ())]; found {
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

func basicModeFunction72Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction72Table[uint(ci.GetJ())]; found {
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

func basicModeFunction7414Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := basicModeFunction7414Table[uint(ci.GetA())]; found {
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

//	--------------------------------------------------------------------------------------------------------------------

func extendedModeFunction05Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction05Table[uint(ci.GetA())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction07Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction07Table[uint(ci.GetJ())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction33Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction33Table[uint(ci.GetJ())]; found {
		return inst(e)
	} else {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode)
	}
}

func extendedModeFunction50Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction50Table[uint(ci.GetA())]; found {
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

func extendedModeFunction72Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction72Table[uint(ci.GetJ())]; found {
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

func extendedModeFunction7414Handler(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	ci := e.GetCurrentInstruction()
	if inst, found := extendedModeFunction7414Table[uint(ci.GetA())]; found {
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
