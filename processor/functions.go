package processor

// FunctionTable maps the basic mode flag to either the basic mode or extended mode function table
var FunctionTable = map[bool]map[uint]func(*Processor) (completed bool, interrupt Interrupt){
	true:  BasicModeFunctionTable,
	false: ExtendedModeFunctionTable,
}

var BasicModeFunctionTable = map[uint]func(*Processor) (completed bool, interrupt Interrupt){
	10: LoadAccumulator,
}

var ExtendedModeFunctionTable = map[uint]func(*Processor) (completed bool, interrupt Interrupt){
	10: LoadAccumulator,
}

// LoadAccumulator loads the content of U under j-field control, and stores it in A(a)
func LoadAccumulator(p *Processor) (completed bool, interrupt Interrupt) {
	completed = false
	interrupt = nil
	// TODO
	//	long operand = getOperand(true, true, true, true);
	//	setExecOrUserARegister((int) _currentInstruction.getA(), operand);
	return
}
