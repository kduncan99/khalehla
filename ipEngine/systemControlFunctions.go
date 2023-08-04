// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

//	TODO StoreProcessorIdentification (SPID)
//	TODO InstructionProcessorControl (IPC)
//	TODO SystemControl (SYSC)

// InitiateAutoRecovery (IAR)
//
//	In a departure from the architecture guide, we *do* allow this in basic mode as well as extended.
//	This is mainly for unit test purposes.
func InitiateAutoRecovery(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	if e.activityStatePacket.GetDesignatorRegister().GetProcessorPrivilege() > 0 {
		return false, pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
	}

	operand, interrupt := e.GetImmediateOperand()
	if interrupt != nil {
		return false, interrupt
	}
	e.Stop(InitiateAutoRecoveryStop, pkg.Word36(operand))
	return true, nil
}
