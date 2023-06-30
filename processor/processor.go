package processor

import (
	"kalehla"
	"kalehla/types"
	"sync"
)

// stop reasons
const (
	InitialStop uint = iota
	ClearedStop
	DebugStop
	DevelopmentStop
	BreakpointStop
	HaltJumpExecutedStop
	ICSBaseRegisterInvalidStop
	ICSOverflowStop
	InitiateAutoRecoveryStop
	L0BaseRegisterInvalidStop
	PanelHaltStop

	//  Interrupt Handler initiated stops...

	InterruptHandlerHardwareFailureStop
	InterruptHandlerOffsetOutOfRangeStop
	InterruptHandlerInvalidBankTypeStop
	InterruptHandlerInvalidLevelBDIStop
)

const L0BDTBaseRegister = 16
const ICSBaseRegister = 26
const ICSIndexRegister = EX1
const RCSBaseRegister = 25
const RCSIndexRegister = EX0

const JumpHistoryTableThreshold = 120 //	Raise interrupt when this many entries exist
const JumpHistoryTableSize = 128      //	Size of the table

type Processor struct {
	mainStorage *kalehla.MainStorage // must be set externally

	activeBaseTable     [16]*ActiveBaseTableEntry //	[0] is unused
	activityStatePacket ActivityStatePacket
	baseRegisters       [32]*BaseRegister
	generalRegisterSet  GeneralRegisterSet

	//	If not nil, describes an interrupt which needs to be handled as soon as possible
	pendingInterrupt Interrupt

	//	See 2.4.2
	//	Should this be saved off somewhere during an interrupt?
	jumpHistory                 [JumpHistoryTableThreshold]types.Word36
	jumpHistoryIndex            int
	jumpHistoryThresholdReached bool

	//	If true, the current (or most recent) instructionType has set the PAR.PC the way it wants,
	//	and we should not increment it for the next instructionType
	preventPCUpdate bool

	mutex         sync.Mutex
	stopReason    uint
	stopDetail    types.Word36
	terminated    bool
	terminateFlag bool
	//	TODO breakpoint settings
}

var storageLocks = map[uint64]*Processor{}
var storageLocksMutex sync.Mutex

// Order of base register selection for Basic Mode address resolution
// when the Basic Mode Base Register Selection Designator Register bit is false
var baseRegisterCandidatesFalse = []uint{12, 14, 13, 15}

// Order of base register selection for Basic Mode address resolution
// when the Basic Mode Base Register Selection Designator Register bit is true
var baseRegisterCandidatesTrue = []uint{13, 15, 12, 14}

// And the composite of the above...
var baseRegisterCandidates = map[bool][]uint{
	false: baseRegisterCandidatesFalse,
	true:  baseRegisterCandidatesTrue,
}

//	external stuffs ----------------------------------------------------------------------------------------------------

// FindBasicModeBank takes a relative address and determines which (if any) of the basic mode banks
// currently based on BDR12-15 is to be selected for that address.
// Set updatedDB31 true if you want the code to update designator register Bit31 in the event we cross
// primary/secondary bank pairs.
// Returns the bank descriptor index (from 12 to 15) for the proper bank descriptor.
// Returns zero if the address is not within any of the based bank limits.
func (p *Processor) FindBasicModeBank(relativeAddress uint, updateDB31 bool) uint {
	db31 := p.activityStatePacket.designatorRegister.BasicModeBaseRegisterSelection
	for tx := 0; tx < 4; tx++ {
		//  See IP PRM 4.4.5 - select the base register from the selection table.
		//  If the bank is void, skip it.
		//  If the program counter is outside the bank limits, skip it.
		//  Otherwise, we found the BDR we want to use.
		brIndex := baseRegisterCandidates[db31][tx]
		bReg := p.baseRegisters[brIndex]
		if p.isWithinLimits(bReg, relativeAddress) {
			if updateDB31 && (tx >= 2) {
				//  address is found in a secondary bank, so we need to flip DB31
				p.activityStatePacket.designatorRegister.BasicModeBaseRegisterSelection = !db31
			}

			return brIndex
		}
	}

	return 0
}

// GetActiveBaseTableEntry retrieves a pointer to the ABET for the indicated base register 0 to 15
func (p *Processor) GetActiveBaseTableEntry(index uint) *ActiveBaseTableEntry {
	return p.activeBaseTable[index]
}

// GetBaseRegister retrieves a pointer to the indicated base register
func (p *Processor) GetBaseRegister(index uint) *BaseRegister {
	return p.baseRegisters[index]
}

// GetExecOrUserRRegisterIndex retrieves the GRS index of either ER{index} or R{index}
// depending upon the setting of designator register ExecRegisterSetSelected
func (p *Processor) GetExecOrUserRRegisterIndex(registerIndex uint) uint {
	if p.activityStatePacket.designatorRegister.ExecRegisterSetSelected {
		return ER0 + registerIndex
	} else {
		return R0 + registerIndex
	}
}

// GetExecOrUserXRegister retrieves a pointer to the index register which corresponds to
// the given register index (0 to 15), and based upon the setting of designator register ExecRegisterSetSelected
func (p *Processor) GetExecOrUserXRegister(registerIndex uint) *IndexRegister {
	index := p.GetExecOrUserXRegisterIndex(registerIndex)
	return (*IndexRegister)(p.generalRegisterSet.GetRegister(index))
}

// GetExecOrUserXRegisterIndex retrieves the GRS index of either EX{index} or X{index}
// depending upon the setting of designator register ExecRegisterSetSelected
func (p *Processor) GetExecOrUserXRegisterIndex(registerIndex uint) uint {
	if p.activityStatePacket.designatorRegister.ExecRegisterSetSelected {
		return EX0 + registerIndex
	} else {
		return X0 + registerIndex
	}
}

// GetGeneralRegisterSet retrieves a pointer to the GRS
func (p *Processor) GetGeneralRegisterSet() *GeneralRegisterSet {
	return &p.generalRegisterSet
}

// GetStopDetail retrieves the detail for the most recent stop
func (p *Processor) GetStopDetail() types.Word36 {
	return p.stopDetail
}

// GetStopReason retrieves the reason code for the most recent stop
func (p *Processor) GetStopReason() uint {
	return p.stopReason
}

// IsStopped indicates whether the processor is stopped
func (p *Processor) IsStopped() bool {
	return p.terminated
}

// PostInterrupt posts a new interrupt, provided that no higher-priority interrupt is already pending.
func (p *Processor) PostInterrupt(i Interrupt) {
	if p.pendingInterrupt != nil {
		if i.GetClass() < p.pendingInterrupt.GetClass() {
			p.pendingInterrupt = i
		}
	}
}

// SetBaseRegister sets the base register identified by brIndex (0 to 15) to the given register
func (p *Processor) SetBaseRegister(brIndex uint, register *BaseRegister) {
	p.baseRegisters[brIndex] = register
}

func (p *Processor) SetExecOrUserRRegister(regIndex uint, value types.Word36) {
	p.generalRegisterSet.SetRegisterValue(p.GetExecOrUserRRegisterIndex(regIndex), value)
}

func (p *Processor) SetExecOrUserXRegister(regIndex uint, value IndexRegister) {
	p.generalRegisterSet.SetRegisterValue(p.GetExecOrUserXRegisterIndex(regIndex), types.Word36(value))
}

// SetProgramCounter sets the program counter in the PAR as well as the preventPCUpdate (aka prevent increment) flag
func (p *Processor) SetProgramCounter(counter uint, preventIncrement bool) {
	p.activityStatePacket.programAddressRegister.SetProgramCounter(counter)
	p.preventPCUpdate = preventIncrement
}

// Start starts the processor, at wherever PAR happens to be set
func (p *Processor) Start() {
	for ix := 0; ix < 32; ix++ {
		p.baseRegisters[ix] = &BaseRegister{}
	}
	for ix := 0; ix < 16; ix++ {
		p.activeBaseTable[ix] = &ActiveBaseTableEntry{}
	}

	p.terminated = false
	p.jumpHistoryIndex = 0
	p.jumpHistoryThresholdReached = false
	p.pendingInterrupt = nil
	go p.run()
}

// Stop stops the processor, while providing a reason and optionally some detail
func (p *Processor) Stop(reason uint, detail types.Word36) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.terminated {
		p.stopReason = reason
		p.stopDetail = detail
		p.terminateFlag = true
	}
}

//	Internal stuffs ----------------------------------------------------------------------------------------------------

// checkAccessibility compares the given key to the lock for this base register, and determines whether
// the requested access (fetch, read, and/or write) are allowed.
// If so, we return true. If not, we post an interrupt and return false.
func (p *Processor) checkAccessibility(
	bReg *BaseRegister,
	fetchFlag bool,
	readFlag bool,
	writeFlag bool,
	accessKey *AccessKey) bool {
	perms := bReg.GetEffectivePermissions(accessKey)
	if p.activityStatePacket.designatorRegister.BasicModeEnabled && fetchFlag && !perms.CanEnter() {
		ssf := uint(040)
		if fetchFlag {
			ssf |= 01
		}
		p.PostInterrupt(NewReferenceViolationInterrupt(ReferenceViolationReadAccess, fetchFlag))
		return false
	} else if readFlag && !perms.CanRead() {
		p.PostInterrupt(NewReferenceViolationInterrupt(ReferenceViolationReadAccess, fetchFlag))
		return false
	} else if writeFlag && !perms.CanWrite() {
		p.PostInterrupt(NewReferenceViolationInterrupt(ReferenceViolationWriteAccess, fetchFlag))
		return false
	} else {
		return true
	}
}

// checkAccessLimits checks the accessibility of a given relative address in the bank described by this
// base register for the given flags, using the given key.
func (p *Processor) checkAccessLimits(
	bReg *BaseRegister,
	relativeAddress uint,
	fetchFlag bool,
	readFlag bool,
	writeFlag bool,
	accessKey *AccessKey) bool {
	if !p.checkAccessLimitsForAddress(bReg, relativeAddress, fetchFlag) {
		return false
	} else if !p.checkAccessibility(bReg, fetchFlag, readFlag, writeFlag, accessKey) {
		return false
	} else {
		return true
	}
}

// checkAccessLimitsForAddress checks whether the relative address is within the limits of the bank
// described by this base register. We only need the fetch flag for posting an interrupt.
// Returns true if access is allowed, posts an interrupt and returns false otherwise
func (p *Processor) checkAccessLimitsForAddress(
	bReg *BaseRegister,
	relativeAddress uint,
	fetchFlag bool) bool {

	// TODO if we try to execute something in GRS - we take ReferenceViolationInterruptClass, 01, 0, 0

	if (relativeAddress < bReg.lowerLimitNormalized) ||
		(relativeAddress > bReg.upperLimitNormalized) {
		p.PostInterrupt(NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, fetchFlag))
		return false
	} else {
		return true
	}
}

// checkAccessLimitsRange checks the access limits for a consecutive range of addresses, starting at the given
// relativeAddress, for the number of addresses. Checks for read and/or write access according to the values given
// for readFlag and writeFlag. Uses the given access key for the determination.
// If we succeed, we return true. If we fail for any reason, we post an interrupt before returning false
func (p *Processor) checkAccessLimitsRange(
	bReg *BaseRegister,
	relativeAddress uint,
	addressCount uint,
	readFlag bool,
	writeFlag bool,
	accessKey *AccessKey) bool {
	if (relativeAddress < bReg.lowerLimitNormalized) ||
		(relativeAddress+addressCount-1 > bReg.upperLimitNormalized) {
		p.PostInterrupt(NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, false))
		return false
	}

	return p.checkAccessibility(bReg, false, readFlag, writeFlag, accessKey)
}

func (p *Processor) clearStorageLocks() {
	storageLocksMutex.Lock()
	for key, value := range storageLocks {
		if value == p {
			delete(storageLocks, key)
		}
	}
	storageLocksMutex.Unlock()
}

// createJumpHistoryTableENtry puts a new entry into the jump history table.
//
//	If we cross the interrupt threshold, set the threshold-reached flag
func (p *Processor) createJumpHistoryTableEntry(absoluteAddress types.Word36) {
	p.jumpHistory[p.jumpHistoryIndex] = absoluteAddress

	if p.jumpHistoryIndex > JumpHistoryTableThreshold {
		p.jumpHistoryThresholdReached = true
	}

	p.jumpHistoryIndex++
	if p.jumpHistoryIndex == JumpHistoryTableSize {
		p.jumpHistoryIndex = 0
	}
}

func (p *Processor) executeCurrentInstruction() {
	//	functions return true if they have completed (normally, or by posting an interrupt).
	//	They return false if they return before completion, but without posting an interrupt.
	//	Generally, a false return results either from a repeated execution instructionType giving
	//	up the processor, or by an indirect basic mode instructionType giving up the processor
	//	before completely developing the operand address.
	p.preventPCUpdate = false
	fTable := FunctionTable[p.activityStatePacket.designatorRegister.BasicModeEnabled]
	if inst, found := fTable[p.activityStatePacket.currentInstruction.GetF()]; found {
		completed, interrupt := inst(p)
		if interrupt != nil {
			p.PostInterrupt(interrupt)
		} else if completed {
			p.clearStorageLocks()
			p.activityStatePacket.indicatorKeyRegister.instructionInF0 = false
			if !p.preventPCUpdate {
				p.activityStatePacket.programAddressRegister.IncrementProgramCounter()
			}
		}
	} else {
		p.PostInterrupt(NewInvalidInstructionInterrupt(InvalidInstructionBadFunctionCode))
	}
}

// fetchInstructionWord retrieves the next instruction word from the appropriate bank.
// For extended mode this is straight-forward.
// For basic mode, we have to hunt around a bit to make sure we pull it from the most appropriate bank.
// If something bad happens, an interrupt is posted and we return false
func (p *Processor) fetchInstructionWord() bool {
	basicMode := p.activityStatePacket.designatorRegister.BasicModeEnabled
	programCounter := p.activityStatePacket.programAddressRegister.GetProgramCounter()

	var bReg *BaseRegister
	if basicMode {
		brIndex := p.FindBasicModeBank(programCounter, true)
		if brIndex == 0 {
			p.PostInterrupt(NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, false))
			return false
		}

		bReg = p.baseRegisters[brIndex]
		if !p.isReadAllowed(bReg) {
			p.PostInterrupt(NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, false))
			return false
		}
	} else {
		bReg = p.baseRegisters[0]
		ok := p.checkAccessLimits(bReg, programCounter, true, false, false, p.activityStatePacket.indicatorKeyRegister.accessKey)
		if !ok {
			return false
		}
	}

	if bReg.voidFlag || bReg.largeSizeFlag {
		p.PostInterrupt(NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, false))
		return false
	}

	pcOffset := uint64(programCounter) - uint64(bReg.lowerLimitNormalized)
	p.activityStatePacket.currentInstruction = InstructionWord(bReg.storage[pcOffset])
	p.activityStatePacket.indicatorKeyRegister.instructionInF0 = true

	return true
}

// findBankDescriptor retrieves a struct to describe the given named bank.
//
//	This is for interrupt handling.
//	The bank name is in L,BDI format.
//	bankLevel level of the bank, 0:7
//	bankDescriptorIndex BDI of the bank 0:077777
func (p *Processor) findBankDescriptor(bankLevel uint, bankDescriptorIndex uint) (*BankDescriptor, bool) {
	// The bank descriptor tables for bank levels 0 through 7 are described by the banks based on B16 through B23.
	// The bank descriptor will be the {n}th bank descriptor in the particular bank descriptor table,
	// where {n} is the bank descriptor index.
	bdRegIndex := bankLevel + 16
	if p.baseRegisters[bdRegIndex].voidFlag {
		p.PostInterrupt(NewAddressingExceptionInterrupt(AddressingExceptionFatal, bankLevel, bankDescriptorIndex))
		return nil, false
	}

	//  bdStorage contains the BDT for the given bank_name level
	//  bdTableOffset indicates the offset into the BDT, where the bank descriptor is to be found.
	bdStorage := p.baseRegisters[bdRegIndex].storage
	bdTableOffset := bankDescriptorIndex + 8
	if bdTableOffset+8 > uint(len(bdStorage)) {
		p.PostInterrupt(NewAddressingExceptionInterrupt(AddressingExceptionFatal, bankLevel, bankDescriptorIndex))
		return nil, false
	}

	//  Create and return a BankDescriptor object
	bd := NewBankDescriptorFromStorage(bdStorage[bdTableOffset : bdTableOffset+8])
	return bd, true
}

func (p *Processor) handleInterrupt() {
	i := p.pendingInterrupt
	p.pendingInterrupt = nil

	// TODO If the Reset Indicator is set and this is a non-initial exigent (non-deferrable) interrupt,
	//   then error halt and set an SCF readable “register” to indicate that a Reset failure occurred.

	//	A hardware interrupt during hardware interrupt handling is a Very Bad Thing
	if i.GetClass() == HardwareCheckInterruptClass &&
		p.activityStatePacket.designatorRegister.FaultHandlingInProgress {
		p.Stop(InterruptHandlerHardwareFailureStop, 0)
		return
	}

	//	Update fields in the ASP
	p.activityStatePacket.indicatorKeyRegister.SetShortStatusField(i.GetShortStatusField())
	p.activityStatePacket.indicatorKeyRegister.SetInterruptClassField(i.GetClass())
	p.activityStatePacket.interruptStatusWord0 = i.GetStatusWord0()
	p.activityStatePacket.interruptStatusWord1 = i.GetStatusWord1()

	//	Make sure the interrupt control stack base register is valid
	if p.baseRegisters[ICSBaseRegister].voidFlag {
		p.Stop(ICSBaseRegisterInvalidStop, 0)
		return
	}

	//	Acquire a stack frame and verify limits
	icsXReg := IndexRegister(p.generalRegisterSet.GetValueOfRegister(ICSIndexRegister))
	icsXReg.DecrementModifier()
	p.generalRegisterSet.SetRegisterValue(ICSIndexRegister, types.Word36(icsXReg))
	stackOffset := icsXReg.GetXM()
	stackFrameSize := icsXReg.GetXI()
	stackFrameLimit := stackOffset + stackFrameSize
	if (stackFrameLimit-1 > p.baseRegisters[ICSBaseRegister].upperLimitNormalized) ||
		(stackOffset < p.baseRegisters[ICSBaseRegister].lowerLimitNormalized) {
		p.Stop(ICSOverflowStop, 0)
		return
	}

	//	Populate the stack frame in memory
	icsStorage := p.baseRegisters[ICSBaseRegister].storage
	if stackFrameLimit >= uint(len(icsStorage)) {
		p.Stop(ICSBaseRegisterInvalidStop, 0)
		return
	}

	icsSlice := p.baseRegisters[ICSBaseRegister].storage[stackOffset:stackFrameLimit]
	icsSlice[0] = p.activityStatePacket.programAddressRegister.GetComposite()
	icsSlice[1] = types.Word36(p.activityStatePacket.designatorRegister.GetComposite())
	icsSlice[2] = p.activityStatePacket.indicatorKeyRegister.GetComposite()
	icsSlice[3] = p.activityStatePacket.quantumTimer
	icsSlice[4] = types.Word36(p.activityStatePacket.currentInstruction)
	icsSlice[5] = i.GetStatusWord0()
	icsSlice[6] = i.GetStatusWord1()
	for ix := 7; ix < int(stackFrameLimit); ix++ {
		icsSlice[ix] = 0
	}

	p.createJumpHistoryTableEntry(p.activityStatePacket.programAddressRegister.GetComposite())
	NewBankManipulatorForInterrupt(p, i).process()
}

func (p *Processor) isReadAllowed(bReg *BaseRegister) bool {
	permissions := bReg.GetEffectivePermissions(p.activityStatePacket.indicatorKeyRegister.accessKey)
	return permissions.CanRead()
}

// isWithinLimits evaluates the given offset within the constraints of the given base register,
//
//	returning true if the offset is within those constraints, else false
func (p *Processor) isWithinLimits(bReg *BaseRegister, offset uint) bool {
	return !bReg.voidFlag &&
		(offset >= bReg.lowerLimitNormalized) &&
		(offset <= bReg.upperLimitNormalized)
}

func (p *Processor) run() {
	for !p.terminateFlag {
		//	Is there a pending interrupt?
		// Are deferrable interrupts allowed?  If not, ignore the interrupt
		if p.pendingInterrupt != nil {
			if !p.pendingInterrupt.IsDeferrable() || p.activityStatePacket.designatorRegister.DeferrableInterruptEnabled {
				p.handleInterrupt()
				//	do not clear pending interrupt here - the interrupt handler may have posted an interrupt
				continue
			}
		}

		//	Are we continuing an interrupted instructionType?
		// See 5.1.3
		//	INF EXRF Action on User Return from interrupt
		//   0   0   Fetch and execute the instructionType addressed by PAR.
		//   1   0   Obtain the instructionType from F0 (rather than using PAR).
		//   1   1   EXR mid-execution. Enter normal EXR logic at the point where the target instructionType has
		//             just been fetched (but not decoded), using F0 as the target instructionType.
		// Note: In the special case where EXR is itself the target of an EX instructionType, mid-execution state will have
		// EXRF clear until the first interrupt point after the EXR instructionType has been fetched.
		if p.activityStatePacket.indicatorKeyRegister.instructionInF0 {
			p.executeCurrentInstruction()
			continue
		}

		//	Fetch an instructionType and execute it
		if p.fetchInstructionWord() {
			p.executeCurrentInstruction()
		}
	}

	storageLocksMutex.Lock()
	defer storageLocksMutex.Unlock()
	for addr, proc := range storageLocks {
		if proc == p {
			delete(storageLocks, addr)
		}
	}

	p.terminated = true
}
