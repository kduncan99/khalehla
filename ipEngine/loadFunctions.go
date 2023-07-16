// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

// LoadAccumulator loads the content of U under j-field control, and stores it in A(a)
func LoadAccumulator(e *InstructionEngine) (completed bool, interrupt Interrupt) {
	operand, interrupt := e.getOperand(true, true, true, true)
	if interrupt != nil {
		return false, interrupt
	}

	e.GetExecOrUserARegister(uint(e.activityStatePacket.currentInstruction.GetA())).SetW(operand)
	return true, nil
}

func LoadIndexRegister(e *InstructionEngine) (completed bool, interrupt Interrupt) {
	operand, interrupt := e.getOperand(true, true, true, true)
	if interrupt != nil {
		return false, interrupt
	}

	e.GetExecOrUserXRegister(uint(e.activityStatePacket.currentInstruction.GetA())).SetW(operand)
	return true, nil
}

func LoadRegister(e *InstructionEngine) (completed bool, interrupt Interrupt) {
	operand, interrupt := e.getOperand(true, true, true, true)
	if interrupt != nil {
		return false, interrupt
	}

	e.GetExecOrUserRRegister(uint(e.activityStatePacket.currentInstruction.GetA())).SetW(operand)
	return true, nil
}
