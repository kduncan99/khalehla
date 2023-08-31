// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

// ExecutiveRequest generates a class 12 signal interrupt, setting ISW0 to U and SSF to 0.
func ExecutiveRequest(e *InstructionEngine) (completed bool) {
	erCode, interrupt := e.GetImmediateOperand()
	if interrupt != nil {
		e.PostInterrupt(interrupt)
		return false
	} else {
		i := pkg.NewSignalInterrupt(pkg.ERSignal, erCode)
		e.PostInterrupt(i)
		return true
	}
}

// SignalCondition generates a class 12 signal interrupt, setting ISW0 to U and SSF to 1.
func SignalCondition(e *InstructionEngine) (completed bool) {
	signalCode, interrupt := e.GetImmediateOperand()
	if interrupt != nil {
		e.PostInterrupt(interrupt)
		return false
	} else {
		i := pkg.NewSignalInterrupt(pkg.SGNLSignal, signalCode)
		e.PostInterrupt(i)
		return true
	}
}

// PreventAllInterruptsAndJump disables deferrable interrupts (by setting DB13) then jumps to U.
// Requires PP==0 for both basic and extended modes.
func PreventAllInterruptsAndJump(e *InstructionEngine) (completed bool) {
	dr := e.GetDesignatorRegister()
	if dr.GetProcessorPrivilege() > 0 {
		i := pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
		e.PostInterrupt(i)
		return false
	}

	op, flip, comp, i := e.GetJumpOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	if comp {
		if flip {
			// the following is in unconditionalJumpFunctions.go
			e.flipDesignatorRegisterBit31()
		}
		dr.SetDeferrableInterruptEnabled(true)
		e.SetProgramCounter(op, true)
	}

	return comp
}

// AllowAllInterruptsAndJump disables deferrable interrupts (by setting DB13) then jumps to U.
// Requires PP==0 only for extended mode.
func AllowAllInterruptsAndJump(e *InstructionEngine) (completed bool) {
	dr := e.GetDesignatorRegister()
	if !dr.IsBasicModeEnabled() && dr.GetProcessorPrivilege() > 0 {
		i := pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadPP)
		e.PostInterrupt(i)
		return false
	}

	op, flip, comp, i := e.GetJumpOperand()
	if i != nil {
		e.PostInterrupt(i)
		return false
	}

	if comp {
		if flip {
			// the following is in unconditionalJumpFunctions.go
			e.flipDesignatorRegisterBit31()
		}
		dr.SetDeferrableInterruptEnabled(false)
		e.SetProgramCounter(op, true)
	}

	return comp
}
