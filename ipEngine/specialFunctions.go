// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

//	TODO DoubleCountBits (DCB)
//	TODO Execute (EX)
//	TODO ExecuteRepeated (EXR)

// NoOperation (NOP) evaluates the HIU field, but takes no other action (it does do x-register incrementation)
func NoOperation(e *InstructionEngine) (completed bool, interrupt pkg.Interrupt) {
	return e.IgnoreOperand()
}
