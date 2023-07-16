// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

type DesignatorRegister struct {
	ActivityLevelQueueMonitorEnabled bool
	FaultHandlingInProgress          bool
	Executive24BitIndexingEnabled    bool
	QuantumTimerEnabled              bool
	DeferrableInterruptEnabled       bool
	processorPrivilege               uint
	BasicModeEnabled                 bool
	ExecRegisterSetSelected          bool
	Carry                            bool
	Overflow                         bool
	CharacteristicUnderflow          bool
	CharacteristicOverflow           bool
	DivideCheck                      bool
	OperationTrapEnabled             bool
	ArithmeticExceptionEnabled       bool
	BasicModeBaseRegisterSelection   bool
	QuarterWordModeEnabled           bool
}

func (dr *DesignatorRegister) GetComposite() uint64 {
	val := uint64(0)
	if dr.ActivityLevelQueueMonitorEnabled {
		val |= 1 << 0
	}
	if dr.FaultHandlingInProgress {
		val |= 1 << 6
	}
	if dr.Executive24BitIndexingEnabled {
		val |= 1 << 11
	}
	if dr.QuantumTimerEnabled {
		val |= 1 << 12
	}
	if dr.DeferrableInterruptEnabled {
		val |= 1 << 13
	}
	val |= uint64(dr.processorPrivilege&0x03) << 14
	if dr.BasicModeEnabled {
		val |= 1 << 16
	}
	if dr.ExecRegisterSetSelected {
		val |= 1 << 17
	}
	if dr.Carry {
		val |= 1 << 18
	}
	if dr.Overflow {
		val |= 1 << 19
	}
	if dr.CharacteristicUnderflow {
		val |= 1 << 21
	}
	if dr.CharacteristicOverflow {
		val |= 1 << 22
	}
	if dr.DivideCheck {
		val |= 1 << 23
	}
	if dr.OperationTrapEnabled {
		val |= 1 << 27
	}
	if dr.ArithmeticExceptionEnabled {
		val |= 1 << 29
	}
	if dr.BasicModeBaseRegisterSelection {
		val |= 1 << 31
	}
	if dr.QuarterWordModeEnabled {
		val |= 1 << 32
	}

	return val
}

var boolTable = []bool{false, true}

func (dr *DesignatorRegister) SetProcessorPrivilege(value uint) {
	dr.processorPrivilege = value & 03
}

func (dr *DesignatorRegister) SetComposite(value uint64) *DesignatorRegister {
	dr.ActivityLevelQueueMonitorEnabled = boolTable[value&01]
	dr.FaultHandlingInProgress = boolTable[(value>>6)&01]
	dr.Executive24BitIndexingEnabled = boolTable[(value>>11)&01]
	dr.QuantumTimerEnabled = boolTable[(value>>12)&01]
	dr.DeferrableInterruptEnabled = boolTable[(value>>13)&01]
	dr.processorPrivilege = uint((value >> 14) & 03)
	dr.BasicModeEnabled = boolTable[(value>>16)&01]
	dr.ExecRegisterSetSelected = boolTable[(value>>17)&01]
	dr.Carry = boolTable[(value>>18)&01]
	dr.Overflow = boolTable[(value>>19)&01]
	dr.CharacteristicUnderflow = boolTable[(value>>21)&01]
	dr.CharacteristicOverflow = boolTable[(value>>22)&01]
	dr.DivideCheck = boolTable[(value>>23)&01]
	dr.OperationTrapEnabled = boolTable[(value>>27)&01]
	dr.ArithmeticExceptionEnabled = boolTable[(value>>29)&01]
	dr.BasicModeBaseRegisterSelection = boolTable[(value>>31)&01]
	dr.QuarterWordModeEnabled = boolTable[(value>>32)&01]

	return dr
}

func NewDesignatorRegisterFromComposite(value uint64) *DesignatorRegister {
	dr := DesignatorRegister{}
	dr.SetComposite(value)
	return &dr
}
