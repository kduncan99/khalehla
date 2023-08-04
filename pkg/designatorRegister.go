// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

var boolTable = []bool{false, true}

type DesignatorRegister struct {
	activityLevelQueueMonitorEnabled bool
	faultHandlingInProgress          bool
	executive24BitIndexingEnabled    bool
	quantumTimerEnabled              bool
	deferrableInterruptEnabled       bool
	processorPrivilege               uint64
	basicModeEnabled                 bool
	execRegisterSetSelected          bool
	carry                            bool
	overflow                         bool
	characteristicUnderflow          bool
	characteristicOverflow           bool
	divideCheck                      bool
	operationTrapEnabled             bool
	arithmeticExceptionEnabled       bool
	basicModeBaseRegisterSelection   bool
	quarterWordModeEnabled           bool
}

func (dr *DesignatorRegister) Clear() {
	dr.activityLevelQueueMonitorEnabled = false
	dr.faultHandlingInProgress = false
	dr.executive24BitIndexingEnabled = false
	dr.quantumTimerEnabled = false
	dr.deferrableInterruptEnabled = false
	dr.processorPrivilege = 0
	dr.basicModeEnabled = false
	dr.execRegisterSetSelected = false
	dr.carry = false
	dr.overflow = false
	dr.characteristicUnderflow = false
	dr.characteristicOverflow = false
	dr.divideCheck = false
	dr.operationTrapEnabled = false
	dr.arithmeticExceptionEnabled = false
	dr.basicModeBaseRegisterSelection = false
	dr.quarterWordModeEnabled = false
}

func (dr *DesignatorRegister) GetBasicModeBaseRegisterSelection() bool {
	return dr.basicModeBaseRegisterSelection
}

func (dr *DesignatorRegister) GetComposite() uint64 {
	val := uint64(0)
	if dr.activityLevelQueueMonitorEnabled {
		val |= 1 << 0
	}
	if dr.faultHandlingInProgress {
		val |= 1 << 6
	}
	if dr.executive24BitIndexingEnabled {
		val |= 1 << 11
	}
	if dr.quantumTimerEnabled {
		val |= 1 << 12
	}
	if dr.deferrableInterruptEnabled {
		val |= 1 << 13
	}
	val |= (dr.processorPrivilege & 0x03) << 14
	if dr.basicModeEnabled {
		val |= 1 << 16
	}
	if dr.execRegisterSetSelected {
		val |= 1 << 17
	}
	if dr.carry {
		val |= 1 << 18
	}
	if dr.overflow {
		val |= 1 << 19
	}
	if dr.characteristicUnderflow {
		val |= 1 << 21
	}
	if dr.characteristicOverflow {
		val |= 1 << 22
	}
	if dr.divideCheck {
		val |= 1 << 23
	}
	if dr.operationTrapEnabled {
		val |= 1 << 27
	}
	if dr.arithmeticExceptionEnabled {
		val |= 1 << 29
	}
	if dr.basicModeBaseRegisterSelection {
		val |= 1 << 31
	}
	if dr.quarterWordModeEnabled {
		val |= 1 << 32
	}

	return val
}

func (dr *DesignatorRegister) GetProcessorPrivilege() uint64 {
	return dr.processorPrivilege
}

func (dr *DesignatorRegister) IsArithmeticExceptionEnabled() bool {
	return dr.arithmeticExceptionEnabled
}

func (dr *DesignatorRegister) IsBasicModeEnabled() bool {
	return dr.basicModeEnabled
}

func (dr *DesignatorRegister) IsCharacteristicOverflowSet() bool {
	return dr.characteristicOverflow
}

func (dr *DesignatorRegister) IsCharacteristicUnderflowSet() bool {
	return dr.characteristicUnderflow
}

func (dr *DesignatorRegister) IsDeferrableInterruptEnabled() bool {
	return dr.deferrableInterruptEnabled
}

func (dr *DesignatorRegister) IsDivideCheckSet() bool {
	return dr.divideCheck
}

func (dr *DesignatorRegister) IsCarrySet() bool {
	return dr.carry
}

func (dr *DesignatorRegister) IsExecRegisterSetSelected() bool {
	return dr.execRegisterSetSelected
}

func (dr *DesignatorRegister) IsExecutive24BitIndexingSet() bool {
	return dr.execRegisterSetSelected
}

func (dr *DesignatorRegister) IsFaultHandlingInProgress() bool {
	return dr.faultHandlingInProgress
}

func (dr *DesignatorRegister) IsOperationTrapEnabled() bool {
	return dr.operationTrapEnabled
}

func (dr *DesignatorRegister) IsOverflowSet() bool {
	return dr.overflow
}

func (dr *DesignatorRegister) IsQuantumTimerEnabled() bool {
	return dr.quantumTimerEnabled
}

func (dr *DesignatorRegister) IsQuarterWordModeEnabled() bool {
	return dr.quarterWordModeEnabled
}

func (dr *DesignatorRegister) SetActivityLevelQueueMonitorEnabled(value bool) *DesignatorRegister {
	dr.activityLevelQueueMonitorEnabled = value
	return dr
}

func (dr *DesignatorRegister) SetArithmeticExceptionEnabled(value bool) *DesignatorRegister {
	dr.arithmeticExceptionEnabled = value
	return dr
}

func (dr *DesignatorRegister) SetBasicModeBaseRegisterSelection(value bool) *DesignatorRegister {
	dr.basicModeBaseRegisterSelection = value
	return dr
}

func (dr *DesignatorRegister) SetBasicModeEnabled(value bool) *DesignatorRegister {
	dr.basicModeEnabled = value
	return dr
}

func (dr *DesignatorRegister) SetComposite(value uint64) *DesignatorRegister {
	dr.activityLevelQueueMonitorEnabled = boolTable[value&01]
	dr.faultHandlingInProgress = boolTable[(value>>6)&01]
	dr.executive24BitIndexingEnabled = boolTable[(value>>11)&01]
	dr.quantumTimerEnabled = boolTable[(value>>12)&01]
	dr.deferrableInterruptEnabled = boolTable[(value>>13)&01]
	dr.processorPrivilege = (value >> 14) & 03
	dr.basicModeEnabled = boolTable[(value>>16)&01]
	dr.execRegisterSetSelected = boolTable[(value>>17)&01]
	dr.carry = boolTable[(value>>18)&01]
	dr.overflow = boolTable[(value>>19)&01]
	dr.characteristicUnderflow = boolTable[(value>>21)&01]
	dr.characteristicOverflow = boolTable[(value>>22)&01]
	dr.divideCheck = boolTable[(value>>23)&01]
	dr.operationTrapEnabled = boolTable[(value>>27)&01]
	dr.arithmeticExceptionEnabled = boolTable[(value>>29)&01]
	dr.basicModeBaseRegisterSelection = boolTable[(value>>31)&01]
	dr.quarterWordModeEnabled = boolTable[(value>>32)&01]

	return dr
}

func (dr *DesignatorRegister) SetDeferrableInterruptEnabled(value bool) *DesignatorRegister {
	dr.deferrableInterruptEnabled = value
	return dr
}

func (dr *DesignatorRegister) SetExecRegisterSetSelected(value bool) *DesignatorRegister {
	dr.execRegisterSetSelected = value
	return dr
}

func (dr *DesignatorRegister) SetExecutive24BitIndexingEnabled(value bool) *DesignatorRegister {
	dr.executive24BitIndexingEnabled = value
	return dr
}

func (dr *DesignatorRegister) SetFaultHandlingInProgress(value bool) *DesignatorRegister {
	dr.faultHandlingInProgress = value
	return dr
}

func (dr *DesignatorRegister) SetQuantumTimerEnabled(value bool) *DesignatorRegister {
	dr.quantumTimerEnabled = value
	return dr
}

func (dr *DesignatorRegister) SetQuarterWordModeEnabled(value bool) *DesignatorRegister {
	dr.quarterWordModeEnabled = value
	return dr
}

func (dr *DesignatorRegister) SetOperationTrapEnabled(value bool) *DesignatorRegister {
	dr.operationTrapEnabled = value
	return dr
}

func (dr *DesignatorRegister) SetProcessorPrivilege(value uint64) *DesignatorRegister {
	dr.processorPrivilege = value & 03
	return dr
}

func NewDesignatorRegisterFromComposite(value uint64) *DesignatorRegister {
	dr := DesignatorRegister{}
	dr.SetComposite(value)
	return &dr
}
