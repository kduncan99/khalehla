// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

// LoadDesignatorRegister (LD) copies the value from U to the DR, excepting those bits which are set-to-zero
func LoadDesignatorRegister(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	if e.activityStatePacket.GetDesignatorRegister().GetProcessorPrivilege() > 0 {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
	}

	result := e.GetOperand(true, true, false, false, false)
	if result.complete && result.interrupt == nil {
		e.activityStatePacket.GetDesignatorRegister().SetComposite(result.operand)
	}

	return result.complete, result.interrupt
}

// StoreDesignatorRegister (SD) stores the content of the DR to the address specified by U
func StoreDesignatorRegister(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	if e.activityStatePacket.GetDesignatorRegister().GetProcessorPrivilege() > 1 {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
	}

	op := e.activityStatePacket.GetDesignatorRegister().GetComposite()
	return e.StoreOperand(false, true, false, false, op)
}

//	TODO LoadProgramControlDesignators (LPD)
//	TODO StoreProgramControlDesignators (SPD)
//	TODO LoadUserDesignators (LUD)
//	TODO StoreUserDesignators (SUD)
//	TODO LoadAddressingEnvironment (LAE) PP==0
//	TODO UserReturn (UR) PP==0
//	TODO AccelerateUserRegisterSet (ACEL) PP<3
//	TODO DecelerateUserRegisterSet (DCEL) PP<3
//	TODO StoreKeyAndQuantumTimer (SKQT) PP<2
//	TODO KeyChange (KCHG) PP==0
