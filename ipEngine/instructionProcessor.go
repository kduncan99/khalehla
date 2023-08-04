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

func NewInstructionProcessor(mainStorage *pkg.MainStorage) *InstructionProcessor {
	return &InstructionProcessor{
		engine: NewEngine(mainStorage),
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

func (p *InstructionProcessor) handleInterrupt() {
	i := p.engine.pendingInterrupt
	p.engine.pendingInterrupt = nil

	// TODO If the Reset Indicator is set and this is a non-initial exigent (non-deferrable) interrupt,
	//   then error halt and set an SCF readable “register” to indicate that a Reset failure occurred.

	//	A hardware interrupt during hardware interrupt handling is a Very Bad Thing
	if i.GetClass() == pkg.HardwareCheckInterruptClass &&
		p.engine.activityStatePacket.GetDesignatorRegister().IsFaultHandlingInProgress() {
		p.engine.Stop(InterruptHandlerHardwareFailureStop, 0)
		return
	}

	if p.engine.IsLoggingInterrupts() {
		fmt.Printf("--{%s}\n", pkg.GetInterruptString(p.engine.pendingInterrupt))
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
	if (stackFrameLimit-1 > br[ICSBaseRegister].GetUpperLimitNormalized()) ||
		(stackOffset < br[ICSBaseRegister].GetLowerLimitNormalized()) {
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

	p.engine.createJumpHistoryTableEntry(asp.GetProgramAddressRegister().GetComposite())
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
		(offset >= bReg.GetLowerLimitNormalized()) &&
		(offset <= bReg.GetUpperLimitNormalized())
}

// run is the coroutine which drives the engine
func (p *InstructionProcessor) run() {
	for !p.terminate {
		//	Is there a pending interrupt?
		// Are deferrable interrupts allowed?  If not, ignore the interrupt
		if p.engine.pendingInterrupt != nil {
			if !p.engine.pendingInterrupt.IsDeferrable() ||
				p.engine.activityStatePacket.GetDesignatorRegister().IsDeferrableInterruptEnabled() {

				p.handleInterrupt()
				//	do not clear pending interrupt here - the interrupt handler may have posted an interrupt
				continue
			}
		}

		// It's okay to do an engine cycle
		p.engine.doCycle()
	}

	p.engine.clearStorageLocks()
	p.isRunning = false
}
