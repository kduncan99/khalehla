// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

//	TODO StoreProcessorIdentification (SPID) PP<3
//	TODO InstructionProcessorControl (IPC) PP==0
//	TODO SystemControl (SYSC) PP==0

// InitiateAutoRecovery (IAR)
//
//	In a departure from the architecture guide, we *do* allow this in basic mode as well as extended.
//	This is mainly for unit test purposes, and may change at any point.
func InitiateAutoRecovery(e *InstructionEngine) (completed bool) {
	if e.activityStatePacket.GetDesignatorRegister().GetProcessorPrivilege() > 0 {
		i := pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
		e.PostInterrupt(i)
		return false
	}

	operand, i := e.GetImmediateOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	e.Stop(InitiateAutoRecoveryStop, pkg.Word36(operand))
	return true
}
