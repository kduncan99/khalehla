// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

//	TODO DoubleCountBits (DCB)
//	TODO Execute (EX)
//	TODO ExecuteRepeated (EXR)

// NoOperation (NOP) evaluates the HIU field, but takes no other action (it does x-register incrementation)
func NoOperation(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	completed, _, interrupt = e.GetJumpOperand(false)
	return
}
