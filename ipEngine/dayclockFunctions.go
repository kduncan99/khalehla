// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
	"time"
)

//	TODO LRD PP==0

// SelectMasterDayclock does nothing - implemented out of a perverse sense of completion.
func SelectMasterDayclock(e *InstructionEngine) (completed bool) {
	if e.GetDesignatorRegister().GetProcessorPrivilege() > 0 {
		i := pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
		e.PostInterrupt(i)
		return false
	}
	return true
}

// ReadMasterDayclock (RMD) transfers the system day-clock, shifted right 5 bits and
// OR'd with a uniqueness counter, to Aa, Aa+1
var lastReportedValue = uint64(0)

func ReadMasterDayclock(e *InstructionEngine) (completed bool) {
	t := uint64(time.Now().UnixMicro()) << 5
	if t == lastReportedValue {
		t++
	}

	lastReportedValue = t
	ci := e.GetCurrentInstruction()
	aReg0 := e.GetExecOrUserARegister(ci.GetA())
	aReg1 := e.GetExecOrUserARegister(ci.GetA() + 1)
	aReg0.SetW(t >> 36)
	aReg1.SetW(t)

	return true
}

//	TODO LMC PP==0
//	TODO SDMN PP==0
//	TODO SDMS PP==0
//	TODO SDMF PP==0
//	TODO RDC PP==0
