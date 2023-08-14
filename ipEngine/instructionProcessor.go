// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"
	"khalehla/pkg"
	"sync"
	"time"
)

// InstructionProcessor implements a (mostly) fully implementation of the 2200 engine architecture.
// Much of the work is handled by the engine which is common to any emulator, but we provide
// the extra wrapping that makes it a fully standalone instruction engine.
// Generally speaking, you run a native 36-bit exec on this.
type InstructionProcessor struct {
	engine    *InstructionEngine
	mutex     sync.Mutex
	isRunning bool
	terminate bool
}

//	external stuffs ----------------------------------------------------------------------------------------------------

func NewInstructionProcessor(name string, mainStorage *pkg.MainStorage, storageLocks *StorageLocks) *InstructionProcessor {
	return &InstructionProcessor{
		engine: NewEngine(name, mainStorage, storageLocks),
	}
}

// GetStopDetail retrieves the detail for the most recent stop
func (p *InstructionProcessor) GetStopDetail() pkg.Word36 {
	return p.engine.stopDetail
}

// GetStopReason retrieves the reason code for the most recent stop
func (p *InstructionProcessor) GetStopReason() StopReason {
	return p.engine.stopReason
}

// IsStopped indicates whether the ipEngine is stopped
func (p *InstructionProcessor) IsStopped() bool {
	return p.engine.stopReason != NotStopped
}

// Start starts the processor goroutine.
// This is not the same as starting up the engine.
func (p *InstructionProcessor) Start() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isRunning {
		return false
	}

	p.isRunning = true // avoid race conditions
	go p.run()
	return true
}

// Stop stops the processor goroutine (presumably just before deleting it)
// This is not the same as stopping the engine itself (which leaves the processor running)
func (p *InstructionProcessor) Stop() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.isRunning {
		return false
	}

	p.terminate = true
	for p.isRunning {
		time.Sleep(time.Second)
	}

	return true
}

//	Internal stuffs ----------------------------------------------------------------------------------------------------

func (p *InstructionProcessor) handleInterrupt(i pkg.Interrupt) {
	//	TODO if we are posting a fault interrupt, back up machine state for the following:
	//  PAR
	//  Indicator/Key_Register
	//  The requirements for the Quantum_Timer are described in 2.2.4
	//  Designator_Register except:
	// – Architecturally_Undefined: The state of DB31 is undefined for all Reference_Violation interrupts caused by a Basic_Mode Jump_to_Address out of limits.
	// – Architecturally_Undefined: DB21, DB22 and DB23 are undefined for Arithmetic_Exception interrupts, and backed up for all other interrupts.
	// – Architecturally_Undefined: DB18 and DB19 are undefined for fault interrupts or mid-execution interrupts on instructions that modify DB18 and DB19.
	//  Xx (and Xa for BT), if index incrementation is specified
	//  Register operands (GRS location(s) specified by the instruction F0.a (F0.ja for JGD))
	//  User or Executive R1 for EXR, BT, BIML, BICL, BIMT, BAO, and the Search and Masked Search instructions. Note: for the case where R1 = 1 with an EXR target, hardware may optionally have INF = 1 and EXRF = 0 in which case save R1 = 0.
	//  Xs and Xd for BIML (see 6.12.6) and BIMT (see 6.12.4); X1 and X2 for BICL (see 6.12.5)
	//  GRS loaded by ACEL need not be backed up. GRS loaded by LRS need not be backed up with
	// the exception of Aa and Xx
	//  For transfers and Base_Register load instructions (except Load Addressing Environment; see 6.19.7), all Base_Registers
	//  For transfer instructions, Executive X0 and, if a Gate is processed, R0 and R1
	//  For transfer and Base_Register load instructions, User X0 and the data in the RCS frame need
	// not be backed up, but the ABT must remain unchanged from the beginning of the instruction
	//  Except as noted above, when an instruction is to load any GRS register and a fault condition is detected on the corresponding source word, the GRS register to be loaded must remain unaltered. Any GRS register written where the corresponding source word causes no fault condition may remain loaded (including any partial-word writes done from valid source words)
	//  Instruction operands (whether GRS or storage operands; see 4.4.2.5 and 4.4.2.6) never need to be backed up, although a model may choose to do so
	//  A Jump_History entry must not be made for an uncompleted instruction.

	// TODO If the Reset Indicator is set and this is a non-initial exigent (non-deferrable) interrupt,
	//   then error halt and set an SCF readable “register” to indicate that a Reset failure occurred.

	//	A hardware interrupt during hardware interrupt handling is a Very Bad Thing
	if i.GetClass() == pkg.HardwareCheckInterruptClass &&
		p.engine.activityStatePacket.GetDesignatorRegister().IsFaultHandlingInProgress() {
		p.engine.Stop(InterruptHandlerHardwareFailureStop, 0)
		return
	}

	if p.engine.IsLoggingInterrupts() {
		fmt.Printf("--{%s}\n", pkg.GetInterruptString(i))
	}

	asp := p.engine.activityStatePacket
	br := p.engine.baseRegisters
	grs := p.engine.generalRegisterSet

	//	Update fields in the ASP
	asp.GetIndicatorKeyRegister().SetShortStatusField(i.GetShortStatusField())
	asp.GetIndicatorKeyRegister().SetInterruptClassField(i.GetClass())
	asp.SetInterruptStatusWord0(i.GetStatusWord0())
	asp.SetInterruptStatusWord1(i.GetStatusWord1())

	//	Make sure the interrupt control stack base register is valid
	if br[ICSBaseRegister].IsVoid() {
		p.engine.Stop(ICSBaseRegisterInvalidStop, 0)
		return
	}

	//	Acquire a stack frame and verify limits
	icsXReg := IndexRegister(grs.GetValueOfRegister(ICSIndexRegister))
	icsXReg.DecrementModifier()
	grs.SetRegisterValue(ICSIndexRegister, pkg.Word36(icsXReg))
	stackOffset := icsXReg.GetXM()
	stackFrameSize := icsXReg.GetXI()
	stackFrameLimit := stackOffset + stackFrameSize
	if (stackFrameLimit-1 > br[ICSBaseRegister].GetBankDescriptor().GetUpperLimitNormalized()) ||
		(stackOffset < br[ICSBaseRegister].GetBankDescriptor().GetLowerLimitNormalized()) {
		p.engine.Stop(ICSOverflowStop, 0)
		return
	}

	//	Populate the stack frame in memory
	icsStorage := br[ICSBaseRegister].GetStorage()
	if stackFrameLimit >= uint64(len(icsStorage)) {
		p.engine.Stop(ICSBaseRegisterInvalidStop, 0)
		return
	}

	icsSlice := br[ICSBaseRegister].GetStorage()[stackOffset:stackFrameLimit]
	icsSlice[0] = asp.GetProgramAddressRegister().GetComposite()
	icsSlice[1] = pkg.Word36(asp.GetDesignatorRegister().GetComposite())
	icsSlice[2] = asp.GetIndicatorKeyRegister().GetComposite()
	icsSlice[3] = asp.GetQuantumTimer()
	icsSlice[4] = pkg.Word36(*asp.GetCurrentInstruction())
	icsSlice[5] = i.GetStatusWord0()
	icsSlice[6] = i.GetStatusWord1()
	for ix := 7; ix < int(stackFrameLimit); ix++ {
		icsSlice[ix] = 0
	}

	// TODO	p.engine.createJumpHistoryTableEntry(asp.GetProgramAddressRegister().GetComposite())
	NewBankManipulatorForInterrupt(p.engine, i).process()
}

func (p *InstructionProcessor) isReadAllowed(bReg *pkg.BaseRegister) bool {
	permissions := bReg.GetEffectivePermissions(p.engine.activityStatePacket.GetIndicatorKeyRegister().GetAccessKey())
	return permissions.CanRead()
}

// isWithinLimits evaluates the given offset within the constraints of the given base register,
//
//	returning true if the offset is within those constraints, else false
func (p *InstructionProcessor) isWithinLimits(bReg *pkg.BaseRegister, offset uint64) bool {
	return !bReg.IsVoid() &&
		(offset >= bReg.GetBankDescriptor().GetLowerLimitNormalized()) &&
		(offset <= bReg.GetBankDescriptor().GetUpperLimitNormalized())
}

// run is the coroutine which drives the engine
func (p *InstructionProcessor) run() {
	for !p.terminate {

		// TODO
		//	Are we continuing an interrupted instruction?
		// See 5.1.3
		//	INF EXRF Action on User Return from interrupt
		//   0   0   Fetch and execute the instruction addressed by PAR.
		//   1   0   Obtain the instruction from F0 (rather than using PAR).
		//   1   1   EXR mid-execution. Enter normal EXR logic at the point where the target instruction has
		//             just been fetched (but not decoded), using F0 as the target instruction.
		// Note: In the special case where EXR is itself the target of an EX instructionT, mid-execution state will have
		// EXRF clear until the first interrupt point after the EXR instruction has been fetched.

		if !p.engine.pendingInterrupts.IsClear() {
			midExec := p.engine.GetInstructionPoint() == MidInstruction
			resolving := p.engine.GetInstructionPoint() == ResolvingAddress
			deferred := !p.engine.GetDesignatorRegister().IsDeferrableInterruptEnabled()
			i := p.engine.pendingInterrupts.Pop(midExec, resolving, deferred)
			if i != nil {
				p.handleInterrupt(i)
			}
		} else {
			// It's okay to do an engine cycle
			p.engine.doCycle()
		}
	}

	p.engine.clearStorageLocks()
	p.isRunning = false
}
