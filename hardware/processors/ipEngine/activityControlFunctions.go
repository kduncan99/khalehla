// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/common"
)

// LoadDesignatorRegister (LD) copies the value from U to the DR, excepting those bits which are set-to-zero
func LoadDesignatorRegister(e *InstructionEngine) (completed bool) {
	if e.activityStatePacket.GetDesignatorRegister().GetProcessorPrivilege() > 0 {
		i := common.NewInvalidInstructionInterrupt(common.InvalidInstructionBadPP)
		e.PostInterrupt(i)
		return false
	}

	result := e.GetOperand(true, true, false, false, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
	} else if result.complete {
		e.activityStatePacket.GetDesignatorRegister().SetComposite(result.operand)
	}

	return result.complete
}

// StoreDesignatorRegister (SD) stores the content of the DR to the address specified by U
func StoreDesignatorRegister(e *InstructionEngine) (completed bool) {
	if e.activityStatePacket.GetDesignatorRegister().GetProcessorPrivilege() > 1 {
		i := common.NewInvalidInstructionInterrupt(common.InvalidInstructionBadPP)
		e.PostInterrupt(i)
		return false
	}

	op := e.activityStatePacket.GetDesignatorRegister().GetComposite()
	comp, i := e.StoreOperand(false, true, false, false, op)
	if i != nil {
		e.PostInterrupt(i)
	}
	return comp
}

// LoadProgramControlDesignators (LPD) loads a subset of the designator register from the immediate value of U.
// U9 -> DB27 (Operation Trap Enable - certain architecture variations actually set this to zero)
// U11 -> DB29 (Arithmetic Exception Enable)
// U12 -> DB30
// U14 -> DB32 (Quarter-Word Selection)
// U15 -> DB33
// U16 -> DB34
// U17 -> DB35
func LoadProgramControlDesignators(e *InstructionEngine) (completed bool) {
	operand, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	currentValue := e.GetDesignatorRegister().GetComposite()
	newValue := (currentValue & 0_777777_777220) | (operand & uint64(0_000557))
	e.GetDesignatorRegister().SetComposite(newValue)

	return true
}

// StoreProgramControlDesignators (SPD)
func StoreProgramControlDesignators(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(false, true, false, false, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
	} else if result.complete {
		valueMasked := result.source.GetW() & 0_777200
		drMasked := e.GetDesignatorRegister().GetComposite() & 0_000577
		newValue := drMasked | valueMasked
		if result.sourceIsGRS {
			result.source.SetW(newValue)
		} else {
			result.source.SetH2(newValue)
		}
	}

	return result.complete
}

// LoadUserDesignators (LUD) is a superset of LPD... loading from (U) instead of U, and
// affecting DB18-19, 21-23, 27, 29-30, and 32-35.
func LoadUserDesignators(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(false, true, false, false, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
	} else if result.complete {
		currentValue := e.GetDesignatorRegister().GetComposite()
		newValue := (currentValue & 0_777777107220) | (result.operand & uint64(0_000000_0670557))
		e.GetDesignatorRegister().SetComposite(newValue)
	}

	return result.complete
}

// StoreUserDesignators (SUD) is a superset of SPD (see LUD above)
func StoreUserDesignators(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(false, true, false, false, false)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
	} else if result.complete {
		newValue := e.GetDesignatorRegister().GetComposite() & 0_000000_777777
		result.source.SetW(newValue)
	}

	return result.complete
}

//	TODO LoadAddressingEnvironment (LAE) PP==0
//	TODO UserReturn (UR) PP==0

// AccelerateUserRegisterSet (ACEL) PP<3 Loads 32 consecutive registers beginning with X0 (or EX0), and the 16
// consecutive registers beginning with R0 (or ER0), from 48 consecutive words beginning at U.
func AccelerateUserRegisterSet(e *InstructionEngine) (completed bool) {
	if e.activityStatePacket.GetDesignatorRegister().GetProcessorPrivilege() > 2 {
		i := common.NewInvalidInstructionInterrupt(common.InvalidInstructionBadPP)
		e.PostInterrupt(i)
		return false
	}

	result := e.GetConsecutiveOperands(false, 48, true)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
	} else if result.complete {
		ux := 0
		grsRegs := e.GetGeneralRegisterSet().GetConsecutiveRegisters(common.X0, 128)

		ix := e.GetExecOrUserXRegisterIndex(common.X0)
		iLimit := ix + 32
		for ix < iLimit {
			grsRegs[ix].SetW(result.source[ux].GetW())
			ix++
			ux++
		}

		ix = e.GetExecOrUserRRegisterIndex(common.R0)
		iLimit = ix + 16
		for ix < iLimit {
			grsRegs[ix].SetW(result.source[ux].GetW())
			ix++
			ux++
		}
	}

	return result.complete
}

// DecelerateUserRegisterSet (DCEL) PP<3 Copies the X, A, and R registers to the 48 words beginning with U.
// (See ACEL).
func DecelerateUserRegisterSet(e *InstructionEngine) (completed bool) {
	if e.activityStatePacket.GetDesignatorRegister().GetProcessorPrivilege() > 2 {
		i := common.NewInvalidInstructionInterrupt(common.InvalidInstructionBadPP)
		e.PostInterrupt(i)
		return false
	}

	result := e.GetConsecutiveOperands(false, 48, true)
	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
	} else if result.complete {
		ux := 0
		grsRegs := e.GetGeneralRegisterSet().GetConsecutiveRegisters(common.X0, 128)

		ix := e.GetExecOrUserXRegisterIndex(common.X0)
		iLimit := ix + 32
		for ix < iLimit {
			result.source[ux].SetW(grsRegs[ix].GetW())
			ix++
			ux++
		}

		ix = e.GetExecOrUserRRegisterIndex(common.R0)
		iLimit = ix + 16
		for ix < iLimit {
			result.source[ux].SetW(grsRegs[ix].GetW())
			ix++
			ux++
		}
	}

	return result.complete
}

//	TODO StoreKeyAndQuantumTimer (SKQT) PP<2
//	TODO KeyChange (KCHG) PP==0
