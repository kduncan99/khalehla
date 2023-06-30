package processor

import "kalehla/types"

type ActivityStatePacket struct {
	//	Virtual address of the current instructionType.
	//	L,BDI of the PAR refers to the bank currently based on B0 (even throughout basic mode execution)
	//	ProgramCounter of the PAR is the relative address of the current or next instructionType to be executed
	//	(relative to B0 in extended mode, or to one of B12-B15 in basic mode)
	programAddressRegister *ProgramAddressRegister

	//	Current operational modes of the processor
	designatorRegister *DesignatorRegister

	//	Interrupt status information, mid-execution control information, pending interrupt indicators,
	//	current access key.
	indicatorKeyRegister *IndicatorKeyRegister

	//	Signed count-down register - preset to the quantum slice value, and decremented per... whatever.
	//	When a negative value is reached with DB12 set, we take a quantum timer interrupt.
	//	In general, this measures the cpu cost for each instructionType executed, which is held cumulative
	//	elsewhere (presumably by the OS). It should not be updated for the UR instructionType.
	quantumTimer types.Word36

	//	Loaded during the fetchInstructionWord() process (our implementation, not architectural)
	//	Updated in one of the following cases:
	//		When an EX or EXR instructionType is executing, the operand is transferred to F0
	//		For indirect addressing, bits 14-35 of the operand are replaced into F0 as the final address
	//			is developed
	currentInstruction InstructionWord

	//	When found in an ICS, these fields contain interrupt status information
	interruptStatusWord0 types.Word36
	interruptStatusWord1 types.Word36
}

// ReadFromBuffer implements the main functionality for the UR instruction
// (see PRM, bank manipulation step 16)
func (asp *ActivityStatePacket) ReadFromBuffer(buffer []types.Word36) {
	asp.programAddressRegister.SetComposite(uint64(buffer[0]))
	asp.designatorRegister.SetComposite(uint64(buffer[1]))
	ssf := asp.indicatorKeyRegister.shortStatusField
	asp.indicatorKeyRegister.SetComposite(uint64(buffer[2])).SetShortStatusField(ssf)
	asp.quantumTimer = buffer[3]
	asp.currentInstruction = InstructionWord(buffer[4])
}

func (asp *ActivityStatePacket) WriteToBuffer(memory []types.Word36) {
	//	TODO
}

// updateCurrentInstruction is for basic mode processing (and possibly for EX/EXR), where we need to
//
//	replace the XHIU fields of F0, leaving FJA fields intact.
func (asp *ActivityStatePacket) updateCurrentInstruction(value *InstructionWord) {
	asp.currentInstruction &= 0777760000000
	asp.currentInstruction |= *value & 017777777
}
